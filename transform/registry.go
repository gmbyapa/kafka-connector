package transforms

import (
	"fmt"
	"github.com/pickme-go/log"
	"strings"
	"sync"
)

type Registry struct {
	transformerMap sync.Map
}

func NewReg() *Registry {
	return new(Registry)
}

var regLogPrefix = "Connector Registry"

func (r *Registry) Init(config map[string]interface{}) []Transformer {

	var transformers []Transformer

	txs := strings.Split(strings.Replace(config[`transforms`].(string), " ", "", -1), `,`)

	// HelloCast
	// transformers : cast1, cast2
	// cast1.p1 : 100
	// cast2.p1 : 300
	for _, tName := range txs {
		if config[fmt.Sprintf(`transforms.%v.type`, tName)] == nil {
			log.Error(log.WithPrefix(castLogPrefix, fmt.Sprintf(`transforms.%v.type not found`, tName)))
			continue
		}
		transType := config[fmt.Sprintf(`transforms.%v.type`, tName)].(string)
		switchTrans := strings.Split(transType, "$")[0]
		switch switchTrans {
		case `Cast`:
			spec := config[fmt.Sprintf(`transforms.%v.spec`, tName)]
			props := make([]CastProps, 0)
			if strings.Contains(spec.(string), ",") {
				specs := strings.Split(spec.(string), ",")
				for _, s := range specs {
					attrType := strings.Split(s, ":")
					props = append(props, CastProps{attrType[0], attrType[1]})
				}

			} else {
				props = append(props, CastProps{"", spec.(string)})
			}

			transformers = append(transformers, Cast{transType, props})
		case `Drop`:
			beha := config[fmt.Sprintf(`transforms.%v.schema.behavior`, tName)]

			transformers = append(transformers, Drop{transType, beha.(string)})
		case `ExtractField`:
			field := config[fmt.Sprintf(`transforms.%v.field`, tName)].(string)

			transformers = append(transformers, ExtractField{transType, field})
		case `ExtractTopic`:
			field := config[fmt.Sprintf(`transforms.%v.field`, tName)].(string)
			var missOrNull bool
			b := config[fmt.Sprintf(`transforms.%v.skip.missing.or.null`, tName)]
			if b != nil {
				missOrNull = b.(bool)
			}

			transformers = append(transformers, ExtractTopic{transType, field, missOrNull})
		case `Flatten`:
			delimiter := config[fmt.Sprintf(`transforms.%v.delimiter`, tName)].(string)

			transformers = append(transformers, Flatten{transType, delimiter})
		case `HoistField`:
			field := config[fmt.Sprintf(`transforms.%v.field`, tName)].(string)

			transformers = append(transformers, HoistField{transType, field})
		case `InsertField`:
			field := config[fmt.Sprintf(`transforms.%v.static.field`, tName)].(string)
			value := config[fmt.Sprintf(`transforms.%v.static.value`, tName)]

			transformers = append(transformers, InsertField{transType, field, value})
		case `MaskField`:
			f := config[fmt.Sprintf(`transforms.%v.fields`, tName)].(string)
			f = strings.Replace(f, " ", "", -1)
			fields := strings.Split(f, ",")

			transformers = append(transformers, MaskField{transType, fields})
		case `ReplaceField`:
			rp := make([]ReplaceFieldProps, 0)
			rn := config[fmt.Sprintf(`transforms.%v.renames`, tName)]
			bl := config[fmt.Sprintf(`transforms.%v.blacklist`, tName)]
			wl := config[fmt.Sprintf(`transforms.%v.whitelist`, tName)]

			var renames, bls, wls []string
			if bl != nil {
				s := strings.Replace(bl.(string), " ", "", -1)
				bls = strings.Split(s, ",")
			}
			if rn != nil {
				rn = strings.Replace(rn.(string), " ", "", -1)
				renames = strings.Split(rn.(string), ",")
			}
			if wl != nil {
				s := strings.Replace(wl.(string), " ", "", -1)
				wls = strings.Split(s, ",")
			}

			for _, rename := range renames {
				props := strings.Split(rename, ":")
				rp = append(rp, ReplaceFieldProps{props[0], props[1]})
			}

			transformers = append(transformers, ReplaceField{transType, bls, wls, rp})
		case `ValueToKey`:
			var s string
			f := config[fmt.Sprintf(`transforms.%v.fields`, tName)]
			if f != nil {
				s = f.(string)
			}
			s = strings.Replace(s, " ", "", -1)
			fields := strings.Split(s, ",")

			transformers = append(transformers, ValueToKey{fields})
		}
	}

	return transformers
}