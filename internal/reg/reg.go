package reg

import (
	"fmt"
	"log"
	"strings"
)

type Reg struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Short       string  `json:"short_description,omitempty"`
	Type        string  `json:"violation_type,omitempty"`
	Vehicle     Vehicle `json:"-"`
	Outdated    bool    `json:"outdated,omitempty"`
}

var either Vehicle = Taxi | FHV

var All []Reg = []Reg{
	// http://www.nyc.gov/html/dot/downloads/pdf/trafrule.pdf
	{Code: "4-12(p)(2)", Description: "no driving in bike lane", Type: "moving", Vehicle: either},
	{Code: "4-08(e)(9)", Description: "no stopping in bike lane", Type: "parking", Vehicle: either},
	{Code: "4-11(c)(6)", Description: "no pickup or discharge of passengers in bike lane", Type: "parking", Vehicle: either, Short: "pickup/discharge in bike lane"},
	{Code: "4-08(a)(4)", Description: "no parking", Type: "parking", Vehicle: either},
	{Code: "4-08(e)(3)", Description: "no parking on sidewalks", Type: "parking", Vehicle: either},
	{Code: "4-08(b)", Description: "no stopping", Type: "parking", Vehicle: either},
	{Code: "4-08(c)", Description: "no standing", Type: "parking", Vehicle: either},
	{Code: "4-08(j)(2)", Description: "obstructed license plate", Type: "parking", Vehicle: either},
	{Code: "4-07(b)(2)", Description: "blocking intersection and crosswalks", Type: "parking", Vehicle: either, Short: "blocking intersection/xwalk"},
	{Code: "4-05(b)(1)", Description: "no u-turns in business district", Type: "moving", Vehicle: either, Short: "no u-turns"},
	{Code: "4-05(a)", Description: "compliance with turning restrictions", Type: "moving", Vehicle: either},
	{Code: "4-12(i)", Description: "no honking in non-danger situations", Type: "parking", Vehicle: either, Short: "no honking"},
	{Code: "4-12(m)", Description: "no driving in bus & right turn only lane", Type: "moving", Vehicle: either, Short: "no driving in bus lane"},
	{Code: "4-04(b)(3)", Description: "no overtaking a vehicle stopped for pedestrians", Type: "moving", Vehicle: either},

	{Code: "NY VTL 1160(a)", Description: "no right from center lane", Type: "moving", Vehicle: either, Short: "no R from center lane"},
	{Code: "NY VTL 1160(b)", Description: "no left from center lane when both two-way streets", Type: "moving", Vehicle: either, Short: "no L from center (@ 2-way)"},
	{Code: "NY VTL 1160(c)", Description: "no left from center lane at one-way street", Type: "moving", Vehicle: either, Short: "no L from center (@ 1-way)"},
	{Code: "NY VTL 1126", Description: "no passing zone", Type: "moving", Vehicle: either},
	{Code: "NY VTL 402(b)", Description: "license plate must not be obstructed", Type: "parking", Vehicle: either, Short: "obstructed license plate", Outdated: true}, // use 4-08(j)(2)
	{Code: "NY VTL 375(12-a)(b)(1)", Description: "no windshield tint below 70%", Type: "other", Vehicle: either, Short: "no tint below 70%"},
	{Code: "NY VTL 375(12-a)(b)(2)", Description: "no side window tint below 70%", Type: "other", Vehicle: either, Short: "no tint below 70%"},
	{Code: "NY VTL 375(30)", Description: "no obstructed view of road", Type: "other", Vehicle: either},
	{Code: "NY VTL 375(1)(b)(i)", Description: "no posters or stickers on windshield", Type: "other", Vehicle: either},
	{Code: "NY VTL 375(12-a)(a)", Description: "no sign in windshield or side windows", Type: "other", Vehicle: either},
	{Code: "NY VTL 1225-c(2)", Description: "cell-phone use while driving", Type: "moving", Vehicle: either},
	{Code: "NY VTL 1203(a)", Description: "park w/in 12 inches of curb (two way street)", Type: "parking", Vehicle: either, Short: "park w/in 12 inches (@ 2-way)"},
	{Code: "NY VTL 1203(b)", Description: "park w/in 12 inches of curb (one way street)", Type: "parking", Vehicle: either, Short: "park w/in 12 inches (@ 1-way)"},
	{Code: "NY VTL 375(41)", Description: "no blue lights except emergency vehicles", Type: "parking", Vehicle: either, Short: "no blue lights"},
	{Code: "NY VTL 1202(a)(1)(a)", Description: "no double parking", Type: "parking", Vehicle: either},
	{Code: "NY VTL 1225-a", Description: "no driving on sidewalks", Type: "moving", Vehicle: either},

	// http://www.nyc.gov/html/tlc/downloads/pdf/rule_book_current_chapter_80.pdf
	// Valid after 10/26/16
	{Code: "80-12(e)", Type: "other", Description: "threats, harassment, abuse", Vehicle: either},
	{Code: "80-12(f)", Type: "other", Description: "use or threat of physical force", Vehicle: either, Short: "use/threat of physical force"},
	{Code: "80-13(a)(3)(ix)", Type: "moving", Description: "yield sign violation", Vehicle: either},
	{Code: "80-13(a)(3)(iii)", Type: "moving", Description: "following too closely (tailgating)", Vehicle: either},
	{Code: "80-13(a)(3)(vi)", Type: "moving", Description: "failing to yield right of way", Vehicle: either, Short: "failing to yield ROW"},
	{Code: "80-13(a)(3)(vii)", Type: "moving", Description: "traffic signal violation", Vehicle: either},
	{Code: "80-13(a)(3)(viii)", Type: "moving", Description: "stop sign violation", Vehicle: either},
	{Code: "80-13(a)(3)(xi)", Type: "moving", Description: "improper passing", Vehicle: either},
	{Code: "80-13(a)(3)(xii)", Type: "moving", Description: "unsafe lane change", Vehicle: either},
	{Code: "80-13(a)(3)(xiii)", Type: "moving", Description: "driving left of center", Vehicle: either},
	{Code: "80-13(a)(3)(xiv)", Type: "moving", Description: "driving in wrong direction", Vehicle: either},
	{Code: "80-13(a)(3)(i)(A)", Type: "moving", Description: "Speeding 1 to 10 mph above speed limit", Short: "Speeding 1-10mph over limit", Vehicle: either},
	{Code: "80-13(a)(3)(i)(B)", Type: "moving", Description: "Speeding 11 to 20 mph above speed limit", Short: "Speeding 11-20mph over limit", Vehicle: either},
	{Code: "80-13(a)(3)(i)(C)", Type: "moving", Description: "Speeding 21 to 30 mph above speed limit", Short: "Speeding 21-30mph over limit", Vehicle: either},
	{Code: "80-13(a)(3)(i)(D)", Type: "moving", Description: "Speeding 31 to 40 mph above speed limit", Short: "Speeding 31-40mph over limit", Vehicle: either},
	{Code: "80-13(a)(3)(i)(E)", Type: "moving", Description: "Speeding 41 or more mph above speed limit", Short: "Speeding >40mph over limit", Vehicle: either},
	{Code: "80-15(b)", Type: "other", Description: "no smoking", Vehicle: either},
	{Code: "80-17(e)(1)(iv)", Type: "other", Description: "must not pickup if not able to accept credit card", Vehicle: either},
	{Code: "80-17(k)(3)", Type: "other", Description: "fare must be calculated by taximeter", Vehicle: either},

	// Valid through 10/25/16
	{Outdated: true, Code: "54-13(a)(3)(ix)", Type: "moving", Description: "yield sign violation", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(ix)", Type: "moving", Description: "yield sign violation", Vehicle: FHV},
	{Outdated: true, Code: "54-13(a)(3)(vi)", Type: "moving", Description: "failing to yield right of way", Vehicle: Taxi, Short: "failing to yield ROW"},
	{Outdated: true, Code: "55-13(a)(3)(vi)", Type: "moving", Description: "failing to yield right of way", Vehicle: FHV, Short: "failing to yield ROW"},
	{Outdated: true, Code: "54-13(a)(3)(vii)", Type: "moving", Description: "traffic signal violation", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(vii)", Type: "moving", Description: "traffic signal violation", Vehicle: FHV},
	{Outdated: true, Code: "54-13(a)(3)(xi)", Type: "moving", Description: "improper passing", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(xi)", Type: "moving", Description: "improper passing", Vehicle: FHV},
	{Outdated: true, Code: "54-13(a)(3)(xii)", Type: "moving", Description: "unsafe lane change", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(xii)", Type: "moving", Description: "unsafe lane change", Vehicle: FHV},
	{Outdated: true, Code: "54-13(a)(3)(xiii)", Type: "moving", Description: "driving left of center", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(xiii)", Type: "moving", Description: "driving left of center", Vehicle: FHV},
	{Outdated: true, Code: "54-13(a)(3)(xiv)", Type: "moving", Description: "driving in wrong direction", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(xiv)", Type: "moving", Description: "driving in wrong direction", Vehicle: FHV},
	{Outdated: true, Code: "54-15(c)", Type: "other", Description: "no smoking", Vehicle: Taxi},
	{Outdated: true, Code: "55-15(c)", Type: "other", Description: "no smoking", Vehicle: FHV},
	{Outdated: true, Code: "54-12(f)", Type: "other", Description: "threats, harassment, abuse", Vehicle: Taxi},
	{Outdated: true, Code: "55-12(e)", Type: "other", Description: "threats, harassment, abuse", Vehicle: FHV},
	{Outdated: true, Code: "54-12(g)", Type: "other", Description: "use or threat of physical force", Vehicle: Taxi, Short: "use/threat of physical force"},
	{Outdated: true, Code: "55-12(f)", Type: "other", Description: "use or threat of physical force", Vehicle: FHV, Short: "use/threat of physical force"},
	{Outdated: true, Code: "54-22(f)", Type: "other", Description: "device must not obstruct view of road", Vehicle: Taxi},
	{Outdated: true, Code: "54-13(a)(3)(i)(A)", Type: "moving", Description: "Speeding 1 to 10 miles above posted speed limit", Vehicle: Taxi},
	{Outdated: true, Code: "54-13(a)(3)(i)(B)", Type: "moving", Description: "Speeding 11 to 20 miles above posted speed limit", Short: "Speeding 11-20mph over limit", Vehicle: Taxi},
	{Outdated: true, Code: "55-13(a)(3)(i)(A)", Type: "moving", Description: "Speeding 1 to 10 miles above posted speed limit", Vehicle: FHV},
	{Outdated: true, Code: "55-13(a)(3)(i)(B)", Type: "moving", Description: "Speeding 11 to 20 miles above posted speed limit", Vehicle: FHV},
}

type Template struct {
	Code        string
	Description string
}

var Templates []Template = []Template{
	{"*", "At <LOCATION> I observed <VEHICLE> <VIOLATION>. Pictures included."},
	{"*", "At <LOCATION> I observed <VEHICLE> <VIOLATION>. Video included."},
	{"4-12(i)", "While riding bike at <LOCATION>, <VEHICLE> tried to intimidate me by honking at me <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While biking at <LOCATION>, I observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While biking at <LOCATION>, I observed <VEHICLE> blocking intersection and causing gridlock including obstructing bike lane <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While trying to cross the street at <LOCATION>, I observed <VEHICLE> blocking crosswalk obstructing pedestrian ROW <VIOLATION>. Pictures included."},
	{"4-07(b)(2)", "While at <LOCATION>, I observed <VEHICLE> blocking intersection and causing gridlock <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, dangerously forcing bikers (including myself) into traffic lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "<VEHICLE> stopped in bike lane, obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-08(e)(9)", "While near <LOCATION> I observed <VEHICLE> stopped in bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "<VEHICLE> was driving in bike lane to avoid waiting in traffic in through lane, obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "While near <LOCATION> I observed <VEHICLE> driving in bike lane to avoid waiting in through lane for other vehicles <VIOLATION>. Pictures included."},
	{"4-12(p)(2)", "While biking on <LOCATION> I observed <VEHICLE> driving in bike lane as a second vehicle lane (it's not) obstructing my use of bike lane <VIOLATION>. Pictures included."},
	{"4-12(m)", "While at <LOCATION> as a pedestrian I observed <VEHICLE> driving in bus only lane (4-7pm M-F) to avoid traffic <VIOLATION>. Pictures included."},
	{"4-12(m)", "While at <LOCATION> as a pedestrian I observed <VEHICLE> driving in bus only lane (4-7pm M-F) to avoid traffic. Driver did not make a right turn or stop to up/discharging passenger at curb <VIOLATION>. Pictures included."},
	{"80-13(a)(3)(vi)", "At <LOCATION>, <VEHICLE> cut me off in the bike lane failing to yield right of way <VIOLATION>. Pictures included."},
	{"80-13(a)(3)(vii)", "At <LOCATION> I observed <VEHICLE> run red light <VIOLATION>. Pictures included. Pictures show light red and vehicle before intersection, and then vehicle proceeding through intersection on red."},
	{"80-13(a)(3)(vii)", "At <LOCATION> I observed <VEHICLE> run red light <VIOLATION>. Video included. Video shows light red and vehicle before intersection, and then vehicle proceeding through intersection on red."},
	{"80-13(a)(3)(xiii)", "At <LOCATION> I observed <VEHICLE> drive left of center yellow line for a block in an effort to avoid traffic <VIOLATION>. Pictures included."},
	{"4-08(j)(2)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of front license plate <VIOLATION>. Pictures included show obstructed view."},
	{"4-08(j)(2)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of front license plate <VIOLATION> a parking violation subject to Commission Rule 80-13(a)(1). Pictures included show obstructed view."},
	{"4-08(j)(2)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of \"T&LC\" text on rear license plate <VIOLATION>. Pictures included show obstructed view."},
	{"4-08(j)(2)", "At <LOCATION> I observed <VEHICLE> with license plate frame obstructing view of \"T&LC\" text on rear license plate <VIOLATION> a parking violation subject to Commission Rule 80-13(a)(1). Pictures included show obstructed view."},
	{"NY VTL 1202(a)(1)(a)", "While biking on <LOCATION> I observed <VEHICLE> double parked (with no driver in vehicle) causing other vehicles to drive in the bike lane <VIOLATION>. Pictures included."},
	{"NY VTL 1160(c)", "While at <LOCATION> I observed <VEHICLE> make a left turn from center lane to avoid turning traffic <LOCATION>. Pictures included."},
}

func FormatTemplate(template, location, vehicle, license, violation string) string {
	switch vehicle {
	case "FHV":
		vehicle = fmt.Sprintf("Respondent Driver of FHV Vehicle with plate %s", license)
	default:
		vehicle = fmt.Sprintf("Respondent Driver of Taxicab Medallion %s", license)
		// Respondent Driver of Street Hail Livery AB544
	}

	template = strings.Replace(template, "<LOCATION>", location, -1)
	template = strings.Replace(template, "<VEHICLE>", vehicle, -1)
	return strings.Replace(template, "<VIOLATION>", fmt.Sprintf("in violation of %s", violation), -1)
}

func (reg Reg) LongCode() string {
	code := reg.Code
	switch {
	case strings.HasPrefix(code, "4-"):
		code = "NYC TR " + code
	case strings.HasPrefix(code, "54-"):
		fallthrough
	case strings.HasPrefix(code, "55-"):
		code = "Commission Rule " + code
	case strings.HasPrefix(code, "80-"):
		code = "Commission Rule " + code
	}
	return code
}

func (reg Reg) String() string {
	code := reg.LongCode()

	switch reg.Type {
	case "", "parking", "moving":
	case ".":
		return code + reg.Description
	default:
		log.Printf("unknown reg type %v", reg.Type)
		return code + reg.Description
	}

	return fmt.Sprintf("%s (%s)", code, reg.Description)
}
