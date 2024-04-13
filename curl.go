package curl

import (
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func CurlToRequest(args string) (*http.Request, error) {
	var (
		url     string
		request string
		data    string
		headers []string
	)
	var rootCmd = &cobra.Command{
		Use:   "curl",
		Short: "A sample command with repeated flags",
		Run: func(cmd *cobra.Command, args []string) {
			url = args[0]
		},
	}

	rootCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "<header/@file> Pass custom header(s) to server")
	rootCmd.Flags().StringVarP(&request, "request", "X", "GET", "<method> Specify request method to use")
	rootCmd.Flags().StringVarP(&data, "data-raw", "", "", "Specify hosts (can be repeated)")
	rootCmd.Flags().StringVarP(&data, "data", "", "d", "Specify hosts (can be repeated)")
	rootCmd.Flags().Bool("location", true, "Follow redirects")
	rootCmd.Flags().Bool("compressed", true, "Follow redirects")

	spArgs := strings.ReplaceAll(args, "\\\n\t", "")
	spArgs = strings.ReplaceAll(spArgs, "\\\n", "")

	argsList := parseArgs(spArgs)

	if strings.HasPrefix(args, "curl") {
		rootCmd.SetArgs(argsList[1:])
	} else {
		rootCmd.SetArgs(argsList)
	}

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	var bodyData io.Reader
	if data != "" {
		bodyData = strings.NewReader(data)
		request = http.MethodPost
	}
	req, err := http.NewRequest(request, url, bodyData)
	if err != nil {
		return nil, err
	}
	for _, header := range headers {
		kvs := strings.Split(header, ":")
		if len(kvs) == 2 {
			req.Header.Add(kvs[0], kvs[1])
		}
	}
	return req, nil
}

func parseArgs(argStr string) []string {
	var parsedArgs []string
	var stack []rune
	var temp string

	for i := 0; i < len(argStr); i++ {
		pos := rune(argStr[i])
		if pos == '\'' || pos == '"' {
			if len(stack) > 0 && pos == stack[len(stack)-1] {
				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, pos)
			}
			temp += string(pos)
		} else if pos == ' ' && len(stack) == 0 {
			parsedArgs = append(parsedArgs, unq(temp))
			temp = ""
		} else {
			temp += string(pos)
		}
	}

	if temp != "" && len(stack) == 0 {
		parsedArgs = append(parsedArgs, unq(temp))
	}

	return parsedArgs
}

func unq(str string) string {
	if str[0] == str[len(str)-1] {
		if strings.HasPrefix(str, "'") || strings.HasPrefix(str, "\"") {
			return str[1 : len(str)-1]
		}
	}
	return str
}
