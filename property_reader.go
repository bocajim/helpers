package helpers

import (
	"os";
	"io";
	"bufio";
	"strconv"
	"strings"
)

// The PropertyReader can be used to read and write property files.
//
// This object can be used to read and write property files.  Below is an example:
// 		var pr PropertyReader
//		err := pr.ReadFile("foo.properties")

var Properties map[string]map[string]string

// Reads a file into a property map.
func LoadPropertyFile(name string, fileName string) (bool,error) {

	if Properties==nil {
		Properties = make(map[string]map[string]string)
	}
	
	properties := make(map[string]string)

	f,e := os.Open(fileName)
	if e!=nil {
		return false,e
	}
	defer f.Close()
	
	rdr := bufio.NewReaderSize(f,128)
	
	var line string
	for ;; {
		line,e = rdr.ReadString('\n')
		if e!=nil && e!=io.EOF {
			return false,e
		}
		if len(line)>0 && !strings.HasPrefix(line,"#") {
			
			line=strings.TrimRight(line,"\r\n")

			kv := strings.Split(line,"=")
			if len(kv)!=2 { continue }
			if len(kv[1])==0 { continue }
			properties[kv[0]] = kv[1]
		}
		if e==io.EOF {
			break
		}
	}
	Properties[name]=properties;
	return true,nil
}

func PropertiesGetString(group string, key string, def string, defp *string) string {
	if v,f:=Properties[group][key];f {
		return v
	}
	if defp==nil {
		return def
	} else {
		return *defp
	}
}

func PropertiesGetInt(group string, key string, def int, defp *int) int {
	if v,f:=Properties[group][key];f {
		if i,e:=strconv.Atoi(v);e==nil {
			return i
		} 
	}
	if defp==nil {
		return def
	} else {
		return *defp
	}
}
