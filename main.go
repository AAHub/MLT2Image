package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	S           = 1024
	FONT_SIZE   = 16.0
	LINE_HEIGHT = 18.0
)

type Mlt struct {
	Name    string `json:"name,omitempty"`
	Path    string `json:"path,omitempty"`
	Updated string `json:"updated,omitempty"`
	Aa      []*Aa  `json:"aa,omitempty"`
}

type Aa struct {
	Value string `json:"value,omitempty"`
}

type Character struct {
	Value   string `json:"value"`
	Code    string `json:"code"`
	Unicode string `json:"unicode"`
}

func main() {
	dirwalk("./input")
}

func dirwalk(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			dirwalk(filepath.Join(dir, file.Name()))
			continue
		}
		path := filepath.Join(dir, file.Name())
		fromFile(file.Name(), path)
	}
}

func fromFile(name string, filePath string) {
	if name == ".DS_Store" {
		return
	}
	if strings.LastIndex(name, ".mlt") == -1 {
		return
	}
	// ファイルを開く
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
		os.Exit(1)
	}

	// 関数return時に閉じる
	defer f.Close()

	// Scannerで読み込む
	lines := make([]string, 0, 300) // ある程度行数が事前に見積もれるようであれば、makeで初期capacityを指定して予めメモリを確保しておくことが望ましい
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}

	var mlt Mlt
	namepos := strings.LastIndex(name, ".")
	new_name := name[:namepos]

	mlt.Name = new_name

	aa := &Aa{}
	for idx, l := range lines {
		l, _ = sjis_to_utf8(l)
		l, _ = escape_html_special_caharacter(l)
		if l == "[SPLIT]" {
			if aa.Value != "" {
				mlt.Aa = append(mlt.Aa, aa)
			}
			aa = &Aa{}
		} else {
			aa.Value = aa.Value + l + "\n"
		}
		if idx == len(lines)-1 {
			if aa.Value != "" {
				mlt.Aa = append(mlt.Aa, aa)
			}
		}
	}

	pos := strings.LastIndex(filePath, ".")
	fileName := filePath[:pos]

	pos = strings.Index(fileName, "/")
	path := fileName[(pos + 1):]
	mlt.Path = path

	for idx, aa := range mlt.Aa {
		fileName = strconv.Itoa(idx)
		lines := strings.Split(aa.Value, "\n")
		outputDirectory := "output/" + mlt.Path
		outputPath := outputDirectory + "/" + fileName + ".png"
		if err := os.MkdirAll(outputDirectory, 0777); err != nil {
			fmt.Println(err)
		}

		if _, err := ConvertTextToImage(lines, outputPath); err != nil {
			fmt.Println(err)
		}
	}
	return
}

func ConvertTextToImage(lines []string, path string) ([]byte, error) {
	// 対象アスキーアートの縦横を図る
	measure := gg.NewContext(S, S)
	if err := measure.LoadFontFace("./Saitamaar.ttf", FONT_SIZE); err != nil {
		return nil, err
	}
	maxWidth := 0.0
	for _, line := range lines {
		w, _ := measure.MeasureString(line)
		if maxWidth <= w {
			maxWidth = w
		}
	}

	// 対象アスキーアートをpngに描画する
	width := int(maxWidth) + 10
	height := int(int(LINE_HEIGHT) * (len(lines) + 1))
	dc := gg.NewContext(width, height)
	if err := dc.LoadFontFace("./Saitamaar.ttf", FONT_SIZE); err != nil {
		return nil, err
	}
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetHexColor("#333333")
	for idx, line := range lines {
		i := float64(idx + 1)
		dc.DrawString(line, 10, LINE_HEIGHT*i)
	}
	dc.Clip()
	dc.SavePNG(path)
	img := dc.Image()

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	ret := buf.Bytes()
	return ret, nil
}

func sjis_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func escape_html_special_caharacter(str string) (string, error) {
	jsonBlob := []byte(`[{"value":"ƒ","code":"&fnof;","unicode":"&#402;"},{"value":"ε","code":"&epsilon;","unicode":"&#917;"},{"value":"κ","code":"&kappa;","unicode":"&#922;"},{"value":"ο","code":"&omicron;","unicode":"&#927;"},{"value":"υ","code":"&upsilon;","unicode":"&#933;"},{"value":"α","code":"&alpha;","unicode":"&#945;"},{"value":"ζ","code":"&zeta;","unicode":"&#950;"},{"value":"λ","code":"&lambda;","unicode":"&#955;"},{"value":"π","code":"&pi;","unicode":"&#960;"},{"value":"υ","code":"&upsilon;","unicode":"&#965;"},{"value":"?","code":"&thetasym;","unicode":"&#977;"},{"value":"′","code":"&prime;","unicode":"&#8242;"},{"value":"ℑ","code":"&image;","unicode":"&#8465;"},{"value":"↑","code":"&uarr;","unicode":"&#8593;"},{"value":"⇐","code":"&larr;","unicode":"&#8656;"},{"value":"∀","code":"&forall;","unicode":"&#8704;"},{"value":"∈","code":"&isin;","unicode":"&#8712;"},{"value":"−","code":"&minus;","unicode":"&#8722;"},{"value":"∠","code":"&ang;","unicode":"&#8736;"},{"value":"∫","code":"&int;","unicode":"&#8747;"},{"value":"≠","code":"&ne;","unicode":"&#8800;"},{"value":"⊃","code":"&sup;","unicode":"&#8835;"},{"value":"⊗","code":"&otimes;","unicode":"&#8855;"},{"value":"?","code":"&lfloor;","unicode":"&#8970;"},{"value":"♠","code":"&spades;","unicode":"&#9824;"},{"value":"α","code":"&alpha;","unicode":"&#913;"},{"value":"ζ","code":"&zeta;","unicode":"&#918;"},{"value":"λ","code":"&lambda;","unicode":"&#923;"},{"value":"π","code":"&pi;","unicode":"&#928;"},{"value":"φ","code":"&phi;","unicode":"&#934;"},{"value":"β","code":"&beta;","unicode":"&#946;"},{"value":"η","code":"&eta;","unicode":"&#951;"},{"value":"μ","code":"&mu;","unicode":"&#956;"},{"value":"ρ","code":"&rho;","unicode":"&#961;"},{"value":"φ","code":"&phi;","unicode":"&#966;"},{"value":"?","code":"&upsih;","unicode":"&#978;"},{"value":"″","code":"&prime;","unicode":"&#8243;"},{"value":"ℜ","code":"&real;","unicode":"&#8476;"},{"value":"→","code":"&rarr;","unicode":"&#8594;"},{"value":"⇑","code":"&uarr;","unicode":"&#8657;"},{"value":"∂","code":"&part;","unicode":"&#8706;"},{"value":"∉","code":"&notin;","unicode":"&#8713;"},{"value":"∗","code":"&lowast;","unicode":"&#8727;"},{"value":"∧","code":"&and;","unicode":"&#8743;"},{"value":"∴","code":"&there4;","unicode":"&#8756;"},{"value":"≡","code":"&equiv;","unicode":"&#8801;"},{"value":"⊄","code":"&nsub;","unicode":"&#8836;"},{"value":"⊥","code":"&perp;","unicode":"&#8869;"},{"value":"?","code":"&rfloor;","unicode":"&#8971;"},{"value":"♣","code":"&clubs;","unicode":"&#9827;"},{"value":"β","code":"&beta;","unicode":"&#914;"},{"value":"η","code":"&eta;","unicode":"&#919;"},{"value":"μ","code":"&mu;","unicode":"&#924;"},{"value":"ρ","code":"&rho;","unicode":"&#929;"},{"value":"χ","code":"&chi;","unicode":"&#935;"},{"value":"γ","code":"&gamma;","unicode":"&#947;"},{"value":"θ","code":"&theta;","unicode":"&#952;"},{"value":"ν","code":"&nu;","unicode":"&#957;"},{"value":"ς","code":"&sigmaf;","unicode":"&#962;"},{"value":"χ","code":"&chi;","unicode":"&#967;"},{"value":"?","code":"&piv;","unicode":"&#982;"},{"value":"‾","code":"&oline;","unicode":"&#8254;"},{"value":"™","code":"&trade;","unicode":"&#8482;"},{"value":"↓","code":"&darr;","unicode":"&#8595;"},{"value":"⇒","code":"&rarr;","unicode":"&#8658;"},{"value":"∃","code":"&exist;","unicode":"&#8707;"},{"value":"∋","code":"&ni;","unicode":"&#8715;"},{"value":"√","code":"&radic;","unicode":"&#8730;"},{"value":"∨","code":"&or;","unicode":"&#8744;"},{"value":"∼","code":"&sim;","unicode":"&#8764;"},{"value":"≤","code":"&le;","unicode":"&#8804;"},{"value":"⊆","code":"&sube;","unicode":"&#8838;"},{"value":"⋅","code":"&sdot;","unicode":"&#8901;"},{"value":"?","code":"&lang;","unicode":"&#9001;"},{"value":"♥","code":"&hearts;","unicode":"&#9829;"},{"value":"γ","code":"&gamma;","unicode":"&#915;"},{"value":"θ","code":"&theta;","unicode":"&#920;"},{"value":"ν","code":"&nu;","unicode":"&#925;"},{"value":"σ","code":"&sigma;","unicode":"&#931;"},{"value":"ψ","code":"&psi;","unicode":"&#936;"},{"value":"δ","code":"&delta;","unicode":"&#948;"},{"value":"ι","code":"&iota;","unicode":"&#953;"},{"value":"ξ","code":"&xi;","unicode":"&#958;"},{"value":"σ","code":"&sigma;","unicode":"&#963;"},{"value":"ψ","code":"&psi;","unicode":"&#968;"},{"value":"•","code":"&bull;","unicode":"&#8226;"},{"value":"⁄","code":"&frasl;","unicode":"&#8260;"},{"value":"ℵ","code":"&alefsym;","unicode":"&#8501;"},{"value":"↔","code":"&harr;","unicode":"&#8596;"},{"value":"⇓","code":"&darr;","unicode":"&#8659;"},{"value":"∅","code":"&empty;","unicode":"&#8709;"},{"value":"∏","code":"&prod;","unicode":"&#8719;"},{"value":"∝","code":"&prop;","unicode":"&#8733;"},{"value":"∩","code":"&cap;","unicode":"&#8745;"},{"value":"∝","code":"&cong;","unicode":"&#8773;"},{"value":"≥","code":"&ge;","unicode":"&#8805;"},{"value":"⊇","code":"&supe;","unicode":"&#8839;"},{"value":"?","code":"&lceil;","unicode":"&#8968;"},{"value":"?","code":"&rang;","unicode":"&#9002;"},{"value":"♦","code":"&diams;","unicode":"&#9830;"},{"value":"δ","code":"&delta;","unicode":"&#916;"},{"value":"ι","code":"&iota;","unicode":"&#921;"},{"value":"ξ","code":"&xi;","unicode":"&#926;"},{"value":"τ","code":"&tau;","unicode":"&#932;"},{"value":"ω","code":"&omega;","unicode":"&#937;"},{"value":"ε","code":"&epsilon;","unicode":"&#949;"},{"value":"κ","code":"&kappa;","unicode":"&#954;"},{"value":"ο","code":"&omicron;","unicode":"&#959;"},{"value":"τ","code":"&tau;","unicode":"&#964;"},{"value":"ω","code":"&omega;","unicode":"&#969;"},{"value":"…","code":"&hellip;","unicode":"&#8230;"},{"value":"℘","code":"&weierp;","unicode":"&#8472;"},{"value":"←","code":"&larr;","unicode":"&#8592;"},{"value":"↵","code":"&crarr;","unicode":"&#8629;"},{"value":"⇔","code":"&harr;","unicode":"&#8660;"},{"value":"∇","code":"&nabla;","unicode":"&#8711;"},{"value":"∑","code":"&sum;","unicode":"&#8721;"},{"value":"∞","code":"&infin;","unicode":"&#8734;"},{"value":"∪","code":"&cup;","unicode":"&#8746;"},{"value":"≈","code":"&asymp;","unicode":"&#8776;"},{"value":"⊂","code":"&sub;","unicode":"&#8834;"},{"value":"⊕","code":"&oplus;","unicode":"&#8853;"},{"value":"?","code":"&rceil;","unicode":"&#8969;"},{"value":"◊","code":"&loz;","unicode":"&#9674;"}]`)
	var characters []Character
	if err := json.Unmarshal(jsonBlob, &characters); err != nil {
		fmt.Println(err)
		return str, err
	}

	for _, character := range characters {
		str = strings.Replace(str, character.Code, character.Value, -1)
	}
	return str, nil
}
