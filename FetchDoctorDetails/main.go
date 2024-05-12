package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/otiai10/gosseract/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PrescriptionImageDescription struct {
	Path              string `json:"path"`
	Data              string `json:"data"`
	DescriptionDetail []byte `json:"description_detail"`
}

var (
	db  *gorm.DB
	err error
)

type DBConfig struct {
	Username  string
	Password  string
	Host      string
	Port      string
	DBName    string
	Charset   string
	ParseTime bool
	Loc       string
}

func connectDB() (*gorm.DB, error) {
	dsn := DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_DBNAME"),
		Charset:  os.Getenv("DB_CHARSET"),
	}

	dsnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.DBName, dsn.Charset)
	var err error
	db, err = gorm.Open(mysql.Open(dsnString), &gorm.Config{})

	if err == nil {
		fmt.Println("Connected Successfully!!")
		return db, err
	}

	return db, err
}

type LeadPrescriptionImage struct {
	Id     int `gorm:"primaryKey"`
	Status int
}

func (PrescriptionImageDescription) TableName() string {
	return "prescription_image_description"
}

func (LeadPrescriptionImage) TableName() string {
	return "lead_prescription_image"
}

const sampleDegrees = `\b(MBBS|MD|MS|PhD|BDS|DDS|DMD|DPM|DPT|DVM|DC|DO|DNB|BSc|BSN|RN|NP|PA|MD(ayurveda)?|MD(Homoeopathy)?|MS(Homoeopathy)?|MCh|DM|MD(Medicine)?|MD(Dermatology)?|MD(Pathology)?|MD(Pediatrics)?|MD(Anesthesia)?|MD(Radiology)?|MD(Physiology)?|MD(Pharmacology)?|MD(Obstetrics)?|MD(Gynecology)?|MD(Community Medicine)?|MD(General Medicine)?|MD(Forensic Medicine)?|MD(Psychiatry)?|MD(Ophthalmology)?|MD(ENT)?|MD(Surgery)?|MD(Orthopedics)?|MS(General Surgery)?|MS(Orthopedics)?|MS(Obstetrics and Gynecology)?|MS(Ophthalmology)?|MS(ENT)?|MS(Anatomy)?|MS(Physiology)?|MS(Pharmacology)?|MS(Anesthesia)?|MS(Orthodontics)?|MS(Prosthodontics)?|MS(Periodontology)?|MS(Oral and Maxillofacial Surgery)?|MS(Endodontics)?|MS(Pedodontics)?|MS(Conservative Dentistry)?|MS(Oral Pathology)?|MS(Public Health)?|MS(Nursing)?|MDS(Orthodontics)?|MDS(Prosthodontics)?|MDS(Periodontology)?|MDS(Oral and Maxillofacial Surgery)?|MDS(Endodontics)?|MDS(Pedodontics)?|MDS(Conservative Dentistry)?|MDS(Oral Pathology)?|MDS(Community Dentistry)?|MDS(Oral Medicine)?|MDS(Public Health Dentistry)?)\b`

func main() {
	connectDB()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var (
		dataCount int
		loopCount int
		offset    = 0
	)

	dataCount = 20800
	loopCount = dataCount / 100

	for i := 0; i < loopCount; i++ {
		query := fmt.Sprintf("SELECT id,prescription_image FROM lead_prescription_image where status = 0 limit %d,1; ", offset)
		fmt.Println("QUERY::[queryEntries]", query)
		queryRow, queryErr := db.Raw(query).Rows()

		if queryErr != nil {
			fmt.Println("ERROR::[queryEntries]", queryErr.Error())
			return
		}

		for queryRow.Next() {
			image := ""
			var id int

			queryScanErr := queryRow.Scan(&id, &image)

			if queryScanErr != nil {
				fmt.Println("ERROR SCAN::[queryEntries]", queryScanErr.Error())
				return
			}

			updErr := updateStatusForLeadPrescriptionImage(id, 1)
			if updErr != nil {
				fmt.Println("ERROR SCAN::[updateStatusForLeadPrescriptionImage]", updErr.Error())
				return
			}

			if image != "" {
				stringMap := make(map[string]string)

				// Parse JSON
				err := json.Unmarshal([]byte(image), &stringMap)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				for _, value := range stringMap {
					fmt.Println("IMAGE PATH VALUE", value)

					if strings.HasSuffix(value, ".jpg") || strings.HasSuffix(value, ".png") || strings.HasSuffix(value, ".jpeg") {
						parsedURL, err := url.Parse(value)
						if err != nil {
							fmt.Println("Error:", err)
							return
						}

						bucket := "s3healthians"
						key := strings.TrimLeft(parsedURL.Path, "/s3healthians/")

						signedURL, err := generateS3SignedURL(bucket, key)
						if err != nil {
							fmt.Println("Error:", err)
							continue
						}

						fmt.Println("SIGNED URL", signedURL)

						err3 := downloadImage(signedURL, "images/image.jpg")

						if err3 != nil {
							fmt.Println("Error:", err3)
							continue
						}

						text, err := extractTextFromImage("images/image.jpg")
						if err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println("Extracted Text:", text)

						cityName := fetchCitiesFromDb(text)
						fmt.Println("Extracted City Name : ", cityName)

						doctorNames := strings.Join(removeDuplicates(extractDoctorNames(text)), ", ")
						fmt.Println("Extracted Doctor Names:", doctorNames)

						state := extractState(text)
						fmt.Println("Extracted State:", state)

						qualifications := FetchDoctorQualifications(text)
						qualificationResult := strings.Join(qualifications, ", ")
						fmt.Println("Extracted Degrees : ", removeDuplicates(qualifications))

						clinicName := strings.Join((extractWordsBeforeClinicAndHospital(text)), ", ")
						fmt.Println("Extracted Clinic : ", clinicName)

						phoneNumber := extractPhoneNumbers(text)
						fmt.Println("Extracted Phone Number : ", phoneNumber)

						detailedDescription := map[string]interface{}{
							"doctorNames":         doctorNames,
							"qualificationResult": qualificationResult,
							"clinicName":          clinicName,
							"phoneNumber":         phoneNumber,
							"cityName":            cityName,
							"state":               state,
						}

						newText, txtErr := json.Marshal(text)
						if txtErr != nil {
							fmt.Println("Error:", txtErr)
							continue
						}

						detailedDescriptionJson, txtErr := json.Marshal(detailedDescription)
						if txtErr != nil {
							fmt.Println("Error:", txtErr)
							continue
						}

						PrescriptionImageDescription := &PrescriptionImageDescription{
							Path:              value,
							Data:              string(newText),
							DescriptionDetail: detailedDescriptionJson,
						}

						insrtErr := insertIntoPrescriptionImageDescription(*PrescriptionImageDescription)

						if insrtErr != nil {
							fmt.Println("ERROR SCAN::[insertIntoPrescriptionImageDescription]", insrtErr.Error())
							return
						}

						updErr := updateStatusForLeadPrescriptionImage(id, 2)
						if updErr != nil {
							fmt.Println("ERROR SCAN::[updateStatusForLeadPrescriptionImage]", updErr.Error())
							return
						}
					} else {
						fmt.Println("Not an image")
						updErr := updateStatusForLeadPrescriptionImage(id, 3)
						if updErr != nil {
							fmt.Println("ERROR SCAN::[updateStatusForLeadPrescriptionImage]", updErr.Error())
							return
						}
					}
				}
			}
		}
		offset = offset + 100
	}
}

func extractTextFromImage(imagePath string) (string, error) {
	tes := ""
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("eng")

	err := client.SetTessdataPrefix("/usr/share/tesseract-ocr/4.00/tessdata/")

	if err != nil {
		fmt.Println(err)
		return tes, err
	}

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return tes, fmt.Errorf("file doesnt exist")
	}
	err1 := client.SetImage(imagePath)
	if err1 != nil {
		return tes, err1
	}

	text, err2 := client.Text()

	if err2 != nil {
		return tes, err2
	}

	return text, err2
}

func fetchCitiesFromDb(text string) string {

	var cities []string
	if err == nil {
		query := "SELECT city_name FROM deal_city"
		rows, err := db.Raw(query).Rows()
		if err != nil {
			return "Queryu Failed[::fetchCitiesFromDb]"
		}
		defer rows.Close()

		for rows.Next() {
			var cityName string
			if err := rows.Scan(&cityName); err != nil {
				return "Queryu Failed[::fetchCitiesFromDb]"
			}
			cities = append(cities, cityName)
		}

		if err := rows.Err(); err != nil {
			return "Queryu Failed[::fetchCitiesFromDb]"
		}
	}

	return findCity(text, cities)
}

func findCity(text string, cities []string) string {
	lowerText := strings.ToLower(text)
	for _, city := range cities {
		lowerCity := strings.ToLower(city)
		if strings.Contains(lowerText, lowerCity) {
			return city
		}
	}
	return ""
}

func extractPhoneNumbers(text string) string {
	re := regexp.MustCompile(`\b\d{3}[-.\s]?\d{3}[-.\s]?\d{4}\b|\b\d{11}\b`)
	match := re.FindString(text)
	return match
}

func extractDoctorNames(text string) []string {
	re := regexp.MustCompile(`(?:Dr\.|Doctor)\s+[A-Z][a-z]+\s*([A-Z][a-z]*)?`)
	matches := re.FindAllString(text, -1)
	var doctorNames []string
	doctorNames = append(doctorNames, matches...)
	return doctorNames
}

func removeDuplicates(arr []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for _, v := range arr {
		if !encountered[v] {
			result = append(result, v)
			encountered[v] = true
		}
	}
	return result
}

func getStateList(states map[string]string) []string {
	stateList := make([]string, 0, len(states))
	for state := range states {
		stateList = append(stateList, state)
	}
	return stateList
}

func getStateMap() map[string]string {
	return map[string]string{
		"andaman and nicobar islands": "AN",
		"andhra pradesh":              "AP",
		"arunachal pradesh":           "AR",
		"assam":                       "AS",
		"bihar":                       "BR",
		"chandigarh":                  "CH",
		"chhattisgarh":                "CT",
		"dadra and nagar haveli":      "DN",
		"daman and diu":               "DD",
		"delhi":                       "DL",
		"goa":                         "GA",
		"gujarat":                     "GJ",
		"haryana":                     "HR",
		"himachal pradesh":            "HP",
		"jammu and kashmir":           "JK",
		"jharkhand":                   "JH",
		"karnataka":                   "KA",
		"kerala":                      "KL",
		"lakshadweep":                 "LD",
		"madhya pradesh":              "MP",
		"maharashtra":                 "MH",
		"manipur":                     "MN",
		"meghalaya":                   "ML",
		"mizoram":                     "MZ",
		"nagaland":                    "NL",
		"odisha":                      "OR",
		"puducherry":                  "PY",
		"punjab":                      "PB",
		"rajasthan":                   "RJ",
		"sikkim":                      "SK",
		"tamil nadu":                  "TN",
		"telangana":                   "TG",
		"tripura":                     "TR",
		"uttar pradesh":               "UP",
		"uttarakhand":                 "UK",
		"west bengal":                 "WB",
	}
}

func extractState(text string) string {
	text = strings.ToLower(text)
	stateMap := getStateMap()
	re := regexp.MustCompile(`\b(` + strings.Join(getStateList(stateMap), "|") + `)\b`)

	match := re.FindString(text)
	if match != "" {
		for state, abbrev := range stateMap {
			if strings.ToLower(match) == abbrev {
				return state
			}
		}
		return match
	}

	return ""
}

func FetchDoctorQualifications(text string) []string {
	var qualifications []string
	degreeRegex := regexp.MustCompile(sampleDegrees)
	degrees := degreeRegex.FindAllString(text, -1)
	for i := 0; i < len(degrees); i++ {
		qualifications = append(qualifications, degrees[i])
	}
	return qualifications
}

func extractWordsBeforeClinicAndHospital(text string) []string {
	pattern := `(?i)(\b[\w\s]+\b\s+(?:Clinic|Hospital)\b)`
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(text)
	if len(match) >= 2 {
		word := match[1]
		return []string{word}
	}

	return nil
}

func generateS3SignedURL(bucket, key string) (string, error) {
	sessionCh := make(chan *session.Session)
	errorCh := make(chan error)

	go func() {
		sess, err := GetAWSSession()
		if err != nil {
			errorCh <- err
			return
		}
		sessionCh <- sess
	}()

	var sess *session.Session
	select {
	case sess = <-sessionCh:
	case err := <-errorCh:
		return "", err
	}

	svc := s3.New(sess)

	expiration := 15 * time.Minute
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GetAWSSession() (*session.Session, error) {
	awsRegion := os.Getenv("AWSREGION")
	awsAccessKeyID := os.Getenv("AWSACCESSKEY")
	awsSecretAccessKey := os.Getenv("AWSSECRETKEY")
	awsProfile := "default"

	awsToken := ""
	t := strings.ToUpper("false")
	debug := false
	if t == "TRUE" {
		debug = true
	}

	if debug {
		log.Printf("Initiating AWS Seesion with AWS_PROFILE = %s, AWS_REGION = %s, AWS_ACCESS_KEY_ID = %s, AWS_SECRET_ACCESS_KEY = %s", awsProfile, awsRegion, awsAccessKeyID, awsSecretAccessKey)
	} else {
		log.Printf("Initiating AWS Seesion with AWS_PROFILE = %s, AWS_REGION = %s", awsProfile, awsRegion)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, awsToken),
	})

	return sess, err
}

func downloadImage(url, destination string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, status code: %d", response.StatusCode)
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Image downloaded successfully to: %s\n", destination)
	return nil
}

func updateStatusForLeadPrescriptionImage(id int, status int) error {
	if err := db.Model(&LeadPrescriptionImage{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func insertIntoPrescriptionImageDescription(insertData PrescriptionImageDescription) error {
	err := db.Create(&insertData).Error
	if err != nil {
		return err
	}

	var leadID int64
	err = db.Raw("SELECT LAST_INSERT_ID()").Row().Scan(&leadID)
	if err != nil {
		return err
	}

	fmt.Println("Last inserted ID:", leadID)
	return nil
}
