package personpb

import jsonmask "github.com/JoeyLearnsToCode/json-mask"

func (p *Person) NeedMaskFields() map[string]jsonmask.MaskHandler {
	return map[string]jsonmask.MaskHandler{
		".name":           jsonmask.MaskHandlerFunc(jsonmask.Name),
		".telephone":      jsonmask.MaskHandlerFunc(jsonmask.Telephone),
		".id":             jsonmask.MaskHandlerFunc(jsonmask.ToAsterisks),
		".address.street": jsonmask.MaskHandlerFunc(jsonmask.ToAsterisks),
	}
}
