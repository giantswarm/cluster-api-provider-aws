/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package userdata

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
)

const (
	bastionBashScript = `{{.Header}}

BASTION_BOOTSTRAP_FILE=bastion_bootstrap.sh
BASTION_BOOTSTRAP=https://s3.amazonaws.com/aws-quickstart/quickstart-linux-bastion/scripts/bastion_bootstrap.sh

curl -s -o $BASTION_BOOTSTRAP_FILE $BASTION_BOOTSTRAP
chmod +x $BASTION_BOOTSTRAP_FILE

# This gets us far enough in the bastion script to be useful.
apt-get -y update && apt-get -y install python-pip
pip install --upgrade pip &> /dev/null

./$BASTION_BOOTSTRAP_FILE --enable true
`
)

var bastionCloudConfig = `#cloud-config
{{template "files" .WriteFiles}}
`

// BastionInput defines the context to generate a bastion instance user data.
type BastionInput struct {
	baseUserData
}

func NewScript(input *BastionInput) (string, error) {
	input.Header = defaultHeader
	return generate("bastion", bastionBashScript, input)
}

func NewCloudConfig(input *BastionInput) (string, error) {
	return generate("cloudconfig", bastionCloudConfig, input)
}

func NewBastion(input *BastionInput) (string, error) {
	cloudConfig, err := NewCloudConfig(input)
	if err != nil {
		return "", err
	}

	script, err := NewScript(input)
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)

	multi := multipart.NewWriter(b)

	_, err = b.Write([]byte(`Content-Type: multipart/mixed; boundary="` + multi.Boundary() + `"`))
	_, err = b.Write([]byte{'\n'})
	_, err = b.Write([]byte("MIME-Version: 1.0\n"))
	_, err = b.Write([]byte("Number-Attachments: 2\n"))
	_, err = b.Write([]byte{'\n'})

	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", `text/cloud-config; charset="us-ascii"`)
	//h.Set("MIME-Version", "1.0")
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Disposition", `attachment; filename="cloud-config.txt"`)

	w, err := multi.CreatePart(h)
	if err != nil {
		return "", err
	}

	_, err = w.Write([]byte(cloudConfig))
	if err != nil {
		return "", err
	}

	h = make(textproto.MIMEHeader)
	h.Set("Content-Type", `text/x-shellscript; charset="us-ascii"`)
	//h.Set("MIME-Version", "1.0")
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Disposition", `attachment; filename="userdata.txt"`)

	w, err = multi.CreatePart(h)
	if err != nil {
		return "", err
	}

	_, err = w.Write([]byte(script))
	if err != nil {
		return "", err
	}

	err = multi.Close()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
