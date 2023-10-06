package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"image"
// 	"image/color"
// 	"image/draw"
// 	"image/png"
// 	"mime/multipart"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
//
//

// 	 "golang.org/x/image/font"
// 	"github.com/gin-gonic/gin"
// )

// func createSampleImage() ([]byte, error) {
// 	// Define the image dimensions
// 	width := 200
// 	height := 100

// 	// Create a new RGBA image
// 	img := image.NewRGBA(image.Rect(0, 0, width, height))

// 	// Fill the image with a white background
// 	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

// 	// Create a new red color
// 	red := color.RGBA{255, 0, 0, 255}

// 	// Define the font and size for drawing text
// 	fontSize := 16
// 	fnt := &truetype.Font{}
// 	fntSize := float64(fontSize)

// 	// Create a new drawer
// 	d := &font.Drawer{
// 		Dst:  img,
// 		Src:  image.NewUniform(red),
// 		Face: truetype.NewFace(fnt, &truetype.Options{Size: fntSize}),
// 	}

// 	// Draw some text on the image
// 	text := "Sample Text"
// 	d.DrawString(text, fixed.P(20, 50))

// 	// Encode the image as a PNG
// 	var buf bytes.Buffer
// 	err := png.Encode(&buf, img)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil

// }

// func TestCreateHostelEndpoint(t *testing.T) {
// 	// Create a new Gin router instance
// 	imageContents, err := createSampleImage()
// 	if err != nil {
// 		t.Errorf("Error creating sample image: %v", err)
// 	}

// 	r := gin.Default()

// 	// Initialize your router and endpoints
// 	// Create a JSON request body
// 	requestBody := map[string]interface{}{
// 		"name":                  "Sample Hostel",
// 		"university_id":         123,
// 		"address":               "123 University St",
// 		"city":                  "Sample City",
// 		"state":                 "Sample State",
// 		"country":               "Sample Country",
// 		"description":           "A sample hostel",
// 		"number_of_units":       50,
// 		"number_of_occupied_units": 25,
// 		"number_of_bedrooms":    100,
// 		"number_of_bathrooms":   50,
// 		"kitchen":               true,
// 		"floor_space":           500.5,
// 		"amenities": []map[string]interface{}{
// 			{"id": 1, "name": "Wi-Fi"},
// 			{"id": 2, "name": "Laundry"},
// 		},
// 		"hostel_fee": map[string]interface{}{
// 			"total_amount": 500.0,
// 			"payment_plan": "Monthly",
// 			"breakdown": map[string]float64{
// 				"water":      50.0,
// 				"electricity": 50.0,
// 			},
// 		},
// 	}

// 	// Encode the JSON request body
// 	jsonRequestBody, _ := json.Marshal(requestBody)

// 	// Create a multipart writer for the request
// 	body := new(bytes.Buffer)
// 	writer := multipart.NewWriter(body)

// 	// Add the JSON request body as a form field
// 	field, _ := writer.CreateFormField("hostel")
// 	_, _ = field.Write(jsonRequestBody)

// 	// Create a file field for the image (replace with your actual file path)
// 	imageFile, _ := writer.CreateFormFile("hostel_images", "sample.jpg")
// 	_, _ = imageFile.Write(imageContents)

// 	// Close the multipart writer
// 	_ = writer.Close()

// 	// Create a new HTTP request with the multipart form data
// 	req := httptest.NewRequest("POST", "/hostels", body)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	// Create a response recorder to capture the response
// 	w := httptest.NewRecorder()

// 	// Serve the HTTP request to the Gin router
// 	r.ServeHTTP(w, req)

// 	// Check the response status code
// 	if w.Code != http.StatusCreated {
// 		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
// 	}

// }

// func TestMain(m *testing.M) {
// 	exitCode := RunTests(m)
// 	os.Exit(exitCode)
// }
