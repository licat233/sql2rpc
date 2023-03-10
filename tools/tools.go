/*
 * @Author: licat
 * @Date: 2023-01-14 11:12:42
 * @LastEditors: licat
 * @LastEditTime: 2023-02-18 10:38:38
 * @Description: licat233@gmail.com
 */
package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

func PluralizedName(name string) string {
	chip := name[len(name)-1:]
	switch chip {
	case "s":
		return name
	case "y":
		return name[:len(name)-1] + "ies"
	case "_":
		return name[:len(name)-1] + "list"
	default:
		return name + "s"
	}
}

func HasInSlice(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func PickMarkContents2(startMark, endMark string, content []byte) ([][]byte, error) {
	if len(content) == 0 {
		return [][]byte{}, nil
	}
	expr := fmt.Sprintf("%s((?s).*?)%s", startMark, endMark)

	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	listArr := reg.FindAllSubmatch(content, -1)
	byteArr := [][]byte{}
	for _, list := range listArr {
		target := list[len(list)-1]
		if len(target) == 0 {
			continue
		}
		byteArr = append(byteArr, target)
	}
	// fmt.Printf("%så¹éç»æ: %#v \n", expr, byteArr)
	return byteArr, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func MakeDir(filename string) error {
	has, err := PathExists(filename)
	if err != nil {
		return err
	}
	if !has {
		dir := path.Dir(filename)
		has, err = PathExists(dir)
		if err != nil {
			return err
		}
		if !has {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RTCFile(filename string) (content string, f *os.File, err error) {
	//read
	fd, e := os.ReadFile(filename)
	if e != nil && !errors.Is(e, os.ErrNotExist) {
		err = e
		return
	}

	//to string
	content = string(fd)
	content = strings.TrimSpace(content)
	content = strings.Trim(content, "\n")

	//è¯»å | æžç©º | åå»º
	f, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)

	return
}

func GetFilename(filename string) string {
	// è·åæä»¶å
	filename = filepath.Base(filename)

	// è·åæä»¶ç±»å
	extension := filepath.Ext(filename)

	return strings.TrimSuffix(filename, extension)
}

func SetFileType(filepath, filetype string) string {
	fileType := path.Ext(filepath)
	if fileType != filetype {
		filename := filepath[0 : len(filepath)-len(fileType)]
		filepath = fmt.Sprintf("%s%s", filename, filetype)
	}
	return filepath
}

func FileRename(oldFilepath, newname string) string {
	// è·åæä»¶æåšç®åœ
	directory := filepath.Dir(oldFilepath)

	// è·åæä»¶å
	// filename := filepath.Base(oldFilepath)

	// è·åæä»¶ç±»å
	extension := filepath.Ext(oldFilepath)

	// filetype := extension[1:]

	newFilename := fmt.Sprintf("%s%s", newname, extension)
	return path.Join(directory, newFilename)
}

func ToCamel(s string) string {
	return strcase.ToCamel(s)
}

func ToLowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}

func ToSnake(s string) string {
	return strcase.ToSnake(strcase.ToCamel(s))
}

func ExecShell(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// è°çšshellåçº§goçæ¬
func UpgradeCurrentProject(currentVersion, projectInfoURL, projectUrl string) error {
	//è·åææ°çæ¬å·
	version, err := GetLatestReleaseVersion(projectInfoURL)
	if err != nil {
		return fmt.Errorf("\n * Failed to get the latest version numberïŒproject info api url: %s \n   error: %s", projectInfoURL, version)
	}
	//å¯¹æ¯çæ¬å·
	if version == currentVersion {
		fmt.Printf(" version: %s\n The current version is the latest versionïŒno need to upgradeïŒ\n", currentVersion)
		return nil
	}

	//åæ£æ¥goæ¯åŠå­åš
	if _, err = exec.LookPath("go"); err != nil {
		//äžå­åšïŒåæç€ºåå®è£go
		return errors.New("\n * warning: \n   go not exist\n   Please install go first")
	}

	// è¿è¡shellåœä»€ïŒè°çšgo installè¿è¡åçº§
	url := strings.ReplaceAll(projectUrl, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	command := fmt.Sprintf("go install %s@latest", url)
	if out, err := ExecShell(command); err != nil {
		return errors.New(out)
	}

	fmt.Printf("\n Upgrade succeeded: %s -> %s\n", currentVersion, version)
	return nil
}

// è·ågithubé¡¹ç®çææ°releaseçæ¬å·
func GetLatestReleaseVersion(projectInfoURL string) (string, error) {
	command := fmt.Sprintf("wget -qO- -t1 -T2 \"%s\" | grep \"tag_name\" | head -n 1 | awk -F \":\" '{print $2}' | sed 's/\\\"//g;s/,//g;s/ //g'", projectInfoURL)
	out, err := ExecShell(command)
	out = strings.TrimSpace(out)
	return out, err
}

// è·ågitçšæ·å
func GetGitUserName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	username := strings.TrimSpace(string(out))
	return username, nil
}

// è·ågitçšæ·é®ç®±
func GetGitUserEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	username := strings.TrimSpace(string(out))
	return username, nil
}

// è·åç³»ç»çšæ·å
func GetOsUserName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.Username, err
}

// è·ååœåçšæ·å
func GetCurrentUserName() string {
	author, err := GetGitUserName()
	if err != nil || author == "" {
		author, _ = GetOsUserName()
	}
	return author
}

func GetCurrentDirectory() (string, error) {
	//è¿åç»å¯¹è·¯åŸ  filepath.Dir(os.Args[0])å»é€æåäžäžªåçŽ çè·¯åŸ
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	//å°\æ¿æ¢æ/
	return strings.Replace(dir, "\\", "/", -1), nil
}

func GetCurrentDirectoryName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	names := strings.Split(dir, "/")
	return names[len(names)-1], nil
}

func FindFilename(dir string, file string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, fileInfo := range files {
		if strings.EqualFold(fileInfo.Name(), file) {
			return fileInfo.Name(), nil
		}
	}
	return "", nil
}

func FindFile(dir string, file string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, fileInfo := range files {
		if strings.EqualFold(fileInfo.Name(), file) {
			return filepath.Join(dir, fileInfo.Name()), nil
		}
	}
	return "", nil
}
