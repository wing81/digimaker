//Author xc, Created on 2019-08-25 22:51
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"

	_ "github.com/xc/digimaker/sitekit/filters"

	"github.com/gorilla/mux"
	pongo2 "gopkg.in/flosch/pongo2.v2"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}
	filename, err := HandleUploadFile(r, "*")
	result := ""
	if err != nil {
		w.WriteHeader(500)
		result = err.Error()
	} else {
		result = filename
	}
	w.Write([]byte(result))
}

//Upload image, return path or error
func UploadImage(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}
	filename, err := HandleUploadFile(r, ".gif,.jpg,.jpeg,.png")
	result := ""
	if err != nil {
		w.WriteHeader(500)
		result = err.Error()
	} else {
		result = filename
	}
	w.Write([]byte(result))
}

//Handler uploaded file, return filename & error
func HandleUploadFile(r *http.Request, filetype string) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := strings.ToLower(handler.Filename)
	//check if file type is allowed
	fileAllowed := false
	filetypeArr := strings.Split(filetype, ",")
	for _, extension := range filetypeArr {
		if extension == "*" || strings.HasSuffix(filename, extension) {
			fileAllowed = true
			break
		}
	}
	if !fileAllowed {
		return "", errors.New("File format not allowed.")
	}

	tempFolder := util.GetConfig("general", "upload_tempfolder")
	tempFolderAbs := util.VarFolder() + "/" + tempFolder

	//Strip file name
	reg := regexp.MustCompile("[^-A-Za-z0-9_]]")
	filename = reg.ReplaceAllString(filename, "_") //filter out all non word characters

	//Write it to temp folder
	tempFile, err := ioutil.TempFile(tempFolderAbs, "upload-*-"+filename)
	defer tempFile.Close()
	if err != nil {
		return "", err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Chmod(0664)
	tempFile.Write(fileContent)
	pathArr := strings.Split(tempFile.Name(), "/")
	tempFilename := pathArr[len(pathArr)-1]
	return tempFolder + "/" + tempFilename, nil
}

func ExportPDF(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	language := r.FormValue("language")
	if language == "" {
		language = "default"
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(errors.New("id not int"), w)
		return
	}
	querier := handler.Querier()
	content, err := querier.FetchByID(idInt)
	if err != nil {
		HandleError(errors.New("Content not found"), w)
		return
	}
	if !permission.CanRead(userID, content, r.Context()) {
		HandleError(errors.New("Need permission"), w)
		return
	}

	tpl := pongo2.Must(pongo2.FromFile(util.AbsHomePath() + "/templates/pdf/main.html"))
	variables := map[string]interface{}{}
	variables["language"] = language
	variables["content"] = content
	authorID := content.Value("author").(int)
	author, _ := querier.FetchByContentID("user", authorID)
	variables["author"] = author

	data, err2 := tpl.ExecuteBytes(pongo2.Context(variables))
	if err2 != nil {
		HandleError(err2, w)
		return
	}

	modified := content.Value("modified").(int)
	name := util.NameToIdentifier(content.GetName()) + "-" + id + "-" + language + "-" + strconv.Itoa(modified)
	pdfFile, err := HtmlToPDF(string(data), name)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Write([]byte(pdfFile))
}

func HtmlToPDF(html string, name string) (string, error) {
	tempFolder := util.GetConfig("general", "var_folder", "dm")
	targetName := "/pdf/" + name + ".pdf"
	target := tempFolder + targetName

	//if exist already, return it
	if _, err := os.Stat(target); err == nil {
		return targetName, nil
	}
	sourceName := "/pdf/" + name + ".html"

	source := tempFolder + sourceName
	ioutil.WriteFile(source, []byte(html), 0777)
	output, _ := exec.Command("wkhtmltopdf", "--javascript-delay", "1500", "-L", "0px", "-R", "0px", "-T", "95px", "-B", "95.5px", "--print-media-type", "--header-html", tempFolder+"/pdf-assets/header.html", "--footer-html", tempFolder+"/pdf-assets/footer.html", source, target).Output()
	log.Println(output)
	// if err != nil {
	// 	return "", err
	// }
	return targetName, nil
}

func GetAllowedLimitations(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

	params := mux.Vars(r)
	operation := params["operation"]
	operation = strings.ReplaceAll(operation, "_", "/")

	allowedOperations := util.GetConfigArr("permission", "rest_allowed_operations", "dm")
	if !util.Contains(allowedOperations, operation) {
		HandleError(errors.New("Operation not allowed"), w, 403)
		return
	}

	_, limits, err := permission.GetUserAccess(userId, operation, r.Context())
	if err != nil {
		HandleError(err, w)
	} else {
		result, _ := json.Marshal(limits)
		w.Write(result)
	}
}

func init() {
	RegisterRoute("/util/uploadfile", UploadFile)
	RegisterRoute("/util/uploadimage", UploadImage)
	RegisterRoute("/util/limitations/{operation}", GetAllowedLimitations)

	RegisterRoute("/pdf/{id}", ExportPDF)
}
