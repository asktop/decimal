package math

import (
	"strconv"
	"testing"

	"github.com/asktop/decimal"
)

func TestAcos(t *testing.T) {
	eps := new(decimal.Big)
	diff := new(decimal.Big)
	for i, tt := range [...]struct {
		x, r string
	}{
		0: {"-1.00", "3.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117068"},
		1: {"-.9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "3.141592653589793238462643383279502884197169399375091678839320861357328389398966901647249128623363298"},
		2: {"-0.50", "2.094395102393195492308428922186335256131446266250070547316629728205210937524139332418689883561411379"},
		3: {"0", "1.570796326794896619231321691639751442098584699687552910487472296153908203143104499314017412671058534"},
		4: {"0.5", "1.047197551196597746154214461093167628065723133125035273658314864102605468762069666209344941780705689"},
		5: {".9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", "1.414213562373095048801688724209698078569671875376948073176679737990732478462107038850387534327641573E-50"},
		6: {"1.00", "0"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			x, _ := new(decimal.Big).SetString(tt.x)
			r, _ := new(decimal.Big).SetString(tt.r)
			z := decimal.WithPrecision(r.Precision())

			Acos(z, x)
			eps.SetMantScale(1, z.Context.Precision)
			diff.Context.Precision = z.Context.Precision
			if z.Cmp(r) != 0 && diff.Sub(r, z).CmpAbs(eps) > 0 {
				t.Errorf(`#%d: Acos(%s)
wanted: %s
got   : %s
diff  : %s
`, i, x, r, z, diff)
			}
		})
	}
}