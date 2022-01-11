package utility
 
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
 
    uuid "github.com/satori/go.uuid"
    "github.com/sirupsen/logrus"
    "github.com/unrolled/render"
)
 
// Log ...
var Log *logrus.Entry
 
// R ...
var R *render.Render
 
// Target ...
type Target string
 
// Prop ...
type Prop map[string]interface{}
 
const (
    // TargetHeymarketPrefix ...
    TargetHeymarketPrefix = "hm"
    // TargetFacebookPrefix ...
    TargetFacebookPrefix = "fb"
    // TargetLinePrefix ...
    TargetLinePrefix = "line"
    // TargetAbcPrefix ...
    TargetAbcPrefix = "abc"
    // TargetGmbPrefix ...
    TargetGmbPrefix = "gmb"
    // TargetWhatsAppPrefix ...
    TargetWhatsAppPrefix = "whatsapp"
)
 
// SetupService ...
func SetupService(log *logrus.Entry, r *render.Render) {
    Log = log
    R = r
}
 
// ReadAll read and send the body in []byte form...
func ReadAll(body io.ReadCloser) []byte {
    b, e := ioutil.ReadAll(body)
    PrintError(e)
    return b
}
 
// NopCloser returns a ReadCloser with a no-op
// Close method wrapping the provided byte body.
func NopCloser(body []byte) io.ReadCloser {
    return ioutil.NopCloser(bytes.NewBuffer(body))
}
 
// Panic check and panic if needed...
func Panic(e error) {
    if e != nil {
        Log.Panic(e)
    }
}
 
// PrintError ...
func PrintError(e error) {
    if e != nil {
        Log.Println(e)
    }
}
 
// Print ...
func Print(str interface{}) {
    Log.Print(str)
}
 
// Unmarshal ...
func Unmarshal(body []byte, v interface{}) interface{} {
    PrintError(json.Unmarshal(body, v))
    return v
}
 
// Marshal ...
func Marshal(v interface{}) []byte {
    data, err := json.Marshal(v)
    PrintError(err)
    return data
}
 
// ParseForm ..
func ParseForm(r *http.Request) {
    PrintError(r.ParseForm())
}
 
// Int64 ...
func Int64(v string) int64 {
    if IsBlank(v) {
        return 0
    }
    intV, err := strconv.ParseInt(v, 10, 0)
    PrintError(err)
    return intV
}
 
// Int ...
func Int(v string) int {
    if IsBlank(v) {
        return 0
    }
    return int(Int64(v))
}
 
// Float642Int ...
func Float642Int(v interface{}) int {
    if IsBlank(v) {
        return 0
    }
    return int(Float64(v))
}
 
// Float64 ...
func Float64(v interface{}) float64 {
    if IsBlank(v) {
        return 0
    }
    return v.(float64)
}
 
// ToString ...
func ToString(v interface{}) string {
    if IsBlank(v) {
        return ""
    }
    return fmt.Sprintf("%v", v)
}
 
// ToInt ...
func ToInt(v interface{}) int {
    if IsBlank(v) {
        return 0
    }
    return v.(int)
}
 
// CleanPhone ...
func CleanPhone(v interface{}) string {
    if IsBlank(v) {
        return ""
    }
    phone := ToString(v)
    phone = strings.Replace(phone, " ", "", -1)
    phone = strings.Replace(phone, "+", "", -1)
    phone = strings.Replace(phone, "(", "", -1)
    phone = strings.Replace(phone, ")", "", -1)
    phone = strings.Replace(phone, "-", "", -1)
    phone = strings.Replace(phone, ".", "", -1)
    return phone
}
 
// PhoneValid ...
func PhoneValid(v interface{}) bool {
    if IsBlank(v) {
        return false
    }
    if len(ToString(v)) < 10 {
        return false
    }
    return true
}
 
// IsBlank ...
func IsBlank(v interface{}) bool {
    if v == 0 {
        return true
    }
    if v == "" {
        return true
    }
    if v == nil {
        return true
    }
    return false
}
 
// E164Phone ...
func E164Phone(v interface{}) string {
    if IsBlank(v) {
        return ""
    }
    phone := ToString(v)
    hasCountryCode := strings.HasPrefix(phone, "1")
    if len(phone) == 11 && hasCountryCode {
        return phone
    }
    if len(phone) == 10 {
        return "1" + phone
    }
    return phone
}
 
// FormatPhone ... The Phone should have US country code in it. For ex., 19700000987 => +1 (970) 000-0987
func FormatPhone(v string) string {
    // if len(v) > 11 {
    //  return v
    // }
    //  p := Split(v, "")
    //phone := fmt.Sprintf("+1 (%v%v%v) %v%v%v-%v%v%v%v", p[1], p[2], p[3], p[4], p[5], p[6], p[7], p[8], p[9], p[10])
    return v
}
 
// ConvertMap ...
func ConvertMap(v interface{}) map[string]interface{} {
    if IsBlank(v) {
        return nil
    }
    return v.(map[string]interface{})
}
 
// JSON2Map ...
func JSON2Map(v interface{}) map[string]interface{} {
    if IsBlank(v) {
        return nil
    }
    rMsg := v.(json.RawMessage)
    mapp := map[string]interface{}{}
    e := json.Unmarshal(rMsg, &mapp)
    PrintError(e)
    return mapp
}
 
// Split ...
func Split(text string, char string) []string {
    s := strings.Split(text, char)
    return s
}
 
// Trim ...
func Trim(text string, char string) string {
    s := strings.Trim(text, char)
    return s
}
 
// UUID ...
func UUID() string {
    return uuid.NewV4().String()
}
 
// Origin Message Origin
func Origin(target string) string {
    targetArr := strings.Split(target, ":")
    origin := targetArr[0]
    switch origin {
    case TargetHeymarketPrefix:
        return "Heymarket"
    case TargetAbcPrefix:
        return "Apple Business Chat"
    case TargetFacebookPrefix:
        return "Facebook"
    case TargetGmbPrefix:
        return "Google"
    case TargetLinePrefix:
        return "Line"
    case TargetWhatsAppPrefix:
        return "WhatsApp"
    }
    return "SMS"
}
 
// ShowPhone this will allow the message to show the phone number or hide the details if not sent via SMS.
func ShowPhone(target string) bool {
    targetArr := strings.Split(target, ":")
    origin := targetArr[0]
    switch origin {
    case TargetHeymarketPrefix, TargetAbcPrefix, TargetFacebookPrefix, TargetGmbPrefix, TargetLinePrefix:
        return false
    }
    return true
}
 
// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str string, start string, end string) (result string) {
    s := strings.Index(str, start)
    if s == -1 {
        return
    }
    s += len(start)
    e := strings.Index(str[s:], end)
    if e == -1 {
        return
    }
    return str[s : e+s]
}
 
// ShopifyMessage ...
func ShopifyMessage(msg string) string {
    return fmt.Sprintf("[Shopify] %v", msg)
}