// Copyright 2015 Yahoo Inc.
// Licensed under the terms of the Apache version 2.0 license. See LICENSE file for terms.

package main

import (
	"fmt"
	"github.com/ardielle/ardielle-go/rdl"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
)

func javaGenerationHeader(banner string) string {
	return fmt.Sprintf("//\n// This file generated by %s\n//", banner)
}

func javaGenerationPackage(schema *rdl.Schema, ns string) string {
	if ns != "" {
		return ns
	}
	return string(schema.Namespace)
}

func camelSnakeToKebab(name string) string {
	s := strings.Replace(name, "_", "-", -1)
	result := make([]rune, 0)
	wasLower := false
	for _, c := range s {
		if unicode.IsUpper(c) {
			if wasLower {
				result = append(result, '-')
			}
			result = append(result, unicode.ToLower(c))
			wasLower = false
		} else {
			result = append(result, c)
			wasLower = true
		}
	}
	return string(result)
}

func javaGenerationRootPath(schema *rdl.Schema, def string) string {
	if def != "" {
		return def
	}
	if schema.Name != "" {
		n := camelSnakeToKebab(string(schema.Name))
		if schema.Version != nil {
			return fmt.Sprintf("/%s/v%d", n, *schema.Version)
		} else {
			return fmt.Sprintf("/%s", n)
		}
	}
	return "/"
}

func javaGenerationDir(outdir string, schema *rdl.Schema, ns string) (string, error) {
	dir := outdir
	if dir == "" {
		dir = "./src/main/java"
	}
	pack := javaGenerationPackage(schema, ns)
	if pack != "" {
		dir += "/" + strings.Replace(pack, ".", "/", -1)
	}
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
	}
	return dir, err
}

func javaGenerateResourceError(schema *rdl.Schema, writer io.Writer, ns string) error {
	funcMap := template.FuncMap{
		"package": func() string {
			s := javaGenerationPackage(schema, ns)
			if s == "" {
				return s
			}
			return "package " + s + ";\n"
		},
	}
	t := template.Must(template.New("util").Funcs(funcMap).Parse(javaResourceErrorTemplate))
	return t.Execute(writer, schema)
}

const javaResourceErrorTemplate = `{{package}}
public class ResourceError {

    public int code;
    public String message;

    public ResourceError code(int code) {
        this.code = code;
        return this;
    }
    public ResourceError message(String message) {
        this.message = message;
        return this;
    }

    public String toString() {
        return "{code: " + code + ", message: \"" + message + "\"}";
    }

}
`

func javaGenerateResourceException(schema *rdl.Schema, writer io.Writer, ns string) error {
	funcMap := template.FuncMap{
		"package": func() string {
			s := javaGenerationPackage(schema, ns)
			if s == "" {
				return s
			}
			return "package " + s + ";\n"
		},
	}
	t := template.Must(template.New("util").Funcs(funcMap).Parse(javaResourceExceptionTemplate))
	return t.Execute(writer, schema)
}

const javaResourceExceptionTemplate = `{{package}}
public class ResourceException extends RuntimeException {
    public final static int OK = 200;
    public final static int CREATED = 201;
    public final static int ACCEPTED = 202;
    public final static int NO_CONTENT = 204;
    public final static int MOVED_PERMANENTLY = 301;
    public final static int FOUND = 302;
    public final static int SEE_OTHER = 303;
    public final static int NOT_MODIFIED = 304;
    public final static int TEMPORARY_REDIRECT = 307;
    public final static int BAD_REQUEST = 400;
    public final static int UNAUTHORIZED = 401;
    public final static int FORBIDDEN = 403;
    public final static int NOT_FOUND = 404;
    public final static int CONFLICT = 409;
    public final static int GONE = 410;
    public final static int PRECONDITION_FAILED = 412;
    public final static int UNSUPPORTED_MEDIA_TYPE = 415;
    public final static int INTERNAL_SERVER_ERROR = 500;
    public final static int NOT_IMPLEMENTED = 501;

    public final static int SERVICE_UNAVAILABLE = 503;

    public static String codeToString(int code) {
        switch (code) {
        case OK: return "OK";
        case CREATED: return "Created";
        case ACCEPTED: return "Accepted";
        case NO_CONTENT: return "No Content";
        case MOVED_PERMANENTLY: return "Moved Permanently";
        case FOUND: return "Found";
        case SEE_OTHER: return "See Other";
        case NOT_MODIFIED: return "Not Modified";
        case TEMPORARY_REDIRECT: return "Temporary Redirect";
        case BAD_REQUEST: return "Bad Request";
        case UNAUTHORIZED: return "Unauthorized";
        case FORBIDDEN: return "Forbidden";
        case NOT_FOUND: return "Not Found";
        case CONFLICT: return "Conflict";
        case GONE: return "Gone";
        case PRECONDITION_FAILED: return "Precondition Failed";
        case UNSUPPORTED_MEDIA_TYPE: return "Unsupported Media Type";
        case INTERNAL_SERVER_ERROR: return "Internal Server Error";
        case NOT_IMPLEMENTED: return "Not Implemented";
        default: return "" + code;
        }
    }

    int code;
    Object data;

    public ResourceException(int code) {
        this(code, new ResourceError().code(code).message(codeToString(code)));
    }

    public ResourceException(int code, Object data) {
        super("ResourceException (" + code + "): " + data);
        this.code = code;
        this.data = data;
    }

    public int getCode() {
        return code;
    }

    public Object getData() {
        return data;
    }

    public <T> T getData(Class<T> cl) {
        return cl.cast(data);
    }

}
`
