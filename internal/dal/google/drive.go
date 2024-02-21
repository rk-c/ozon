package google

import (
	"bytes"
	"google.golang.org/api/drive/v3"
	"log"
)

func UploadFile(filename string, buf *bytes.Buffer, srv *drive.Service) {

	parent := FindFileByName("SA", srv)
	if parent == nil {
		log.Fatalln("There is no parent file")
	}

	f := &drive.File{Name: filename,
		Parents:  []string{parent.Id},
		MimeType: "text/tsv",
	}

	// Create and upload the file
	_, err := srv.Files.
		Create(f).
		Media(buf). //context.Background(), file, fileInf.Size(), baseMimeType).
		ProgressUpdater(func(now, size int64) {}).
		Do()

	if err != nil {
		log.Fatalln(err)
	}

}

func DeleteFile(filename string, srv *drive.Service) {
	file := FindFileByName(filename, srv)
	if file == nil {
		return
	}
	id := file.Id
	err := srv.Files.Delete(id).Do()

	if err != nil {
		log.Fatalln(err)
	}
}

func FindFileByName(filename string, srv *drive.Service) *drive.File {

	res, err := srv.Files.List().Do()

	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range res.Files {
		if file.Name == filename {
			return file
		}
	}
	log.Printf("Cannot find file %s", filename)
	return nil
}
