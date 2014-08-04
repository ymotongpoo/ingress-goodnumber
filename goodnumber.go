package ingressgoodnumber

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"appengine"
)

const pageTemplate = `<!doctype html>
<html>
  <head>
    <title>Ingress Good Number</title>
	<link rel="stylesheet" type="text/css" href="/main.css" />
    <script type="text/javascript" src="/main.js"></script>
  </head>
  <body>
    <h1>Ingress Good Number</h1>
	<p>"Ingress Good Number" tell you minimum actions to achieve "good number" APs, such as rounded numbers, straight numbers, and repeated numbers.</p>
    <div id='apinput'>
      Your current AP: <input type="text" id="ap" /><input type="button" value="check" id="checkBtn" onclick="goodnumber.postAp()" />
    </div>
    <div id="result"><div>
  </body>
</html>`

const (
	generatorCap = 10
	pattenrsCap  = 20
	maxUint      = ^uint(0)
)

var apGain = []uint{
	1750, // Full deploy
	1563, // Create a CF
	1199, // Destroy a CF
	625,  // Capture a portal
	375,  // Complete a portal
	313,  // Create a link
	262,  // Destroy a link
	125,  // Place a resonator or mod
	100,  // Hack ememy portal
	75,   // Destroy a resonator
	65,   // Upgrade others' resonator
	10,   // Recharge a portal
}

// StatusRequest is a struct defining input data.
type StatusRequest struct {
	AP uint `json:"ap"`
}

// RestActionResponse is a struct defining output data to client.
type RestActionResponse struct {
	Target        uint `json:"target"`
	FullDeploy    uint `json:"full deploy"`
	CreateCF      uint `json:"create control field"`
	DestroyCF     uint `json:"destropy control field"`
	CapturePortal uint `json:"capture portal"`
	CompPortal    uint `json:"complete portal"`
	CreateLink    uint `json:"create link"`
	DestroyLink   uint `json:"destroy link"`
	PlaceRes      uint `json:"place resonator"`
	Hack          uint `json:"hack portal"`
	DestroyRes    uint `json:"destroy resonator"`
	UpgradeRes    uint `json:"upgrade resonator"`
	Recharge      uint `json:"recharge"`
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postHandler(w, r)
	case "GET":
		getHandler(w, r)
	default:
		fmt.Fprintf(w, "This endpoint only support GET or POST methods")
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, pageTemplate)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	c := appengine.NewContext(r)

	var status StatusRequest
	err := decoder.Decode(&status)
	if err != nil {
		message := fmt.Sprintf("An error occured during parsing: %v", err)
		http.Error(w, message, 400)
		c.Errorf("%v", message)
		return
	}
	gn := genGoodNumbers(status.AP)
	c.Infof("AP: %v", status.AP)
	target := <-gn
	pattern := findPattern(status.AP, target)
	action := NewRestActionResponse(target, pattern)
	resp, err := json.Marshal(action)
	if err != nil {
		http.Error(w, err.Error(), 500)
		c.Errorf("%v", err.Error())
		return
	}
	fmt.Fprintf(w, "%v", string(resp))
}

// NewRestActionResponse converts AP list into a struct.
func NewRestActionResponse(target uint, pattern map[uint]uint) *RestActionResponse {
	return &RestActionResponse{
		Target:        target,
		FullDeploy:    pattern[1750],
		CreateCF:      pattern[1563],
		DestroyCF:     pattern[1199],
		CapturePortal: pattern[625],
		CompPortal:    pattern[375],
		CreateLink:    pattern[313],
		PlaceRes:      pattern[125],
		Hack:          pattern[100],
		DestroyRes:    pattern[75],
		UpgradeRes:    pattern[65],
		Recharge:      pattern[10],
	}
}

func genGoodNumbers(ap uint) <-chan uint {
	gn := make(chan uint, generatorCap)
	go func(num uint) {
		digit := numDigits(num)
		round := uint(math.Pow10(digit))
		repdigit := repdigitOf(digit)
		seqdigit := seqdigitOf(digit)

		roundbase := num/round + 1
		repbase := num/repdigit + 1

		nearestRound := round * roundbase
		nearestRep := repdigit * repbase
		var nearestSeq uint
		if seqdigit > num {
			nearestSeq = seqdigit
		} else {
			nearestSeq = seqdigitOf(digit + 1)
		}
		x, y, z := min3(nearestRound, nearestRep, nearestSeq)
		gn <- x
		gn <- y
		gn <- z
	}(ap)
	return gn
}

func findPattern(ap, target uint) map[uint]uint {
	gap := target - ap
	patterns := make([]uint, gap+1)
	track := make([]uint, gap+1)

	// initialize
	for i := uint(0); i < gap+1; i++ {
		patterns[i] = maxUint
	}
	patterns[0] = 0

	// find solution
	for i := uint(0); i < gap+1; i++ {
		min := maxUint
		for _, n := range apGain {
			k := i - n
			if k >= 0 && k < i {
				if patterns[k] < min {
					min = patterns[k]
					patterns[i] = patterns[k] + 1
					track[i] = k
				}
			}
		}
	}

	// find pattern
	if patterns[gap] != maxUint {
		result := createCounterMap()
		for p := gap; ; p = track[p] {
			if track[p] == 0 {
				result[p]++
				break
			}
			result[p-track[p]]++
		}
		return result
	}
	return createCounterMap()
}

// Find order of exponent
func numDigits(num uint) int {
	digit := 0
	for {
		num = num / 10
		if num == 0 {
			break
		}
		digit++
	}
	return digit
}

// repdigitOf returns repdigit with `digit` digits
func repdigitOf(digit int) uint {
	num := uint(1)
	for i := 0; i < digit; i++ {
		num = num*10 + 1
	}
	return num
}

// seqdigitOf returns sequential number with `digit` digits
func seqdigitOf(digit int) uint {
	num := uint(1)
	for i := uint(0); i < uint(digit); i++ {
		num = num*10 + (i+2)%10
	}
	return num
}

// min3 sorts 3 assingments and return them in acsending order
func min3(x, y, z uint) (uint, uint, uint) {
	if x < y {
		if x < z {
			if y < z {
				return x, y, z
			}
			return x, z, y
		}
		return z, x, y
	}
	if z < x {
		if z < y {
			return z, y, x
		}
		return y, z, x
	}
	return y, x, z
}

func createCounterMap() map[uint]uint {
	counter := make(map[uint]uint)
	for _, k := range apGain {
		counter[k] = uint(0)
	}
	return counter
}
