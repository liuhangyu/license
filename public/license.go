package public

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type License struct {
	LicenseUUID string            `json:"licensever,omitempty"`  //license 唯一编号
	ProductName string            `json:"productname,omitempty"` //产品名称
	MachineID   string            `json:"machineid,omitempty"`   //机器ID
	ExpiresAt   int64             `json:"expiresat,omitempty"`   //过期时间
	IssuedAt    int64             `json:"issuedat,omitempty"`    //签发时间
	CustomKV    map[string]string `json:"customkv,omitempty"`
}

func GenerateLicense(uuid string, productName string, machineID string, expires int64, kv map[string]string) *License {
	return &License{
		LicenseUUID: uuid,
		ProductName: productName,
		MachineID:   machineID,
		ExpiresAt:   expires,
		IssuedAt:    time.Now().Unix(),
		CustomKV:    kv,
	}
}

func VerifyLicense(productName string, machine string, licenseBytes []byte, isVerify bool) (*License, error) {
	l := new(License)
	if err := json.Unmarshal(licenseBytes, l); err != nil {
		return nil, err
	}

	if isVerify {
		err := l.Valid(productName, machine)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (c *License) Valid(productName string, machine string) error {
	var vErr error
	now := time.Now().Unix()

	//比较产品名称
	if c.CompareProductName(productName) == false {
		vErr = ErrNoMatchProName
	}

	//比较过期时间
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr = ErrLicenseExpired.SetErrText(fmt.Sprintf("license is expired by %v", delta))
	}

	//比较签发时间
	if c.VerifyIssuedAt(now, false) == false {
		vErr = ErrBeforeIssued
	}

	//比较机器是否与license匹配
	if c.CompareMachine(machine) == false {
		vErr = ErrNoMatchMachineID
	}
	return vErr
}

//已经过期返回0,未过期返回剩余的秒数
func (c *License) GetExpiresAt() int64 {
	now := time.Now().Unix()
	delta := time.Unix(c.ExpiresAt, 0).Sub(time.Unix(now, 0))
	if delta <= 0 { //
		return 0
	}
	return int64(delta.Seconds())
}

func (c *License) GetEndTime() string {
	return time.Unix(c.ExpiresAt, 0).String()
}

//比较产品名
func (c *License) CompareProductName(productName string) bool {
	return strings.Compare(c.ProductName, productName) == 0
}

//比较过期时间
func (c *License) VerifyExpiresAt(now int64, req bool) bool {
	if c.ExpiresAt == 0 {
		return !req
	}
	return now <= c.ExpiresAt
}

//比较签发时间
func (c *License) VerifyIssuedAt(now int64, req bool) bool {
	if c.IssuedAt == 0 {
		return !req
	}
	return now >= c.IssuedAt
}

//比较机器ID
func (c *License) CompareMachine(machineID string) bool {
	return strings.Compare(c.MachineID, machineID) == 0
}

//编码
func (c *License) ToBytes() ([]byte, error) {
	var (
		jsonValue []byte
		err       error
	)
	if jsonValue, err = json.Marshal(c); err != nil {
		return nil, err
	}

	return jsonValue, nil
}

func ToLicense(lb []byte) (*License, error) {
	l := new(License)
	err := json.Unmarshal(lb, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func BytesToLicense(license string) (*License, error) {
	parts := strings.Split(license, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%s", "license contains an invalid number of segments")
	}

	plainBytes, err := DecodeSegment(parts[0])
	if err != nil {
		return nil, err
	}

	l := new(License)
	err = json.Unmarshal(plainBytes, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
