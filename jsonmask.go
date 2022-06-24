package jsonmask

import (
	"encoding/json"
	"fmt"
	"strings"
)

const Sep = "."
const MaxDepth = 10

type NeedMask interface {
	// NeedMaskFields 返回需要掩码处理的字段，key是字段的path（格式： .path.to.filedName ），value是掩码处理器
	NeedMaskFields() map[string]MaskHandler
}

func JsonMaskedQuietly(data interface{}) []byte {
	bytes, _ := JsonMasked(data)
	return bytes
}

// JsonMasked 输出掩码处理的json，入参需要实现 NeedMask 接口，否则输出的是无掩码普通json
func JsonMasked(data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JsonMasked: %w", err)
	}

	if needMask, ok := data.(NeedMask); ok {
		needMaskFields := needMask.NeedMaskFields()

		dataMap := make(map[string]interface{})
		err = json.Unmarshal(jsonBytes, &dataMap)
		if err != nil {
			return nil, fmt.Errorf("JsonMasked: %w", err)
		}

		maskedMap := MaskMapOrSlice(needMaskFields, dataMap)
		jsonBytes, err = json.Marshal(maskedMap)
		if err != nil {
			return nil, fmt.Errorf("JsonMasked: %w", err)
		}
	}
	return jsonBytes, nil
}

// MaskMapOrSlice 对map、slice进行掩码处理，原位更新，返回map或slice。
//  @param needMaskFields 见 NeedMask.NeedMaskFields 注释
func MaskMapOrSlice(needMaskFields map[string]MaskHandler, rawData interface{}) (maskedData interface{}) {
	return dfsMask(needMaskFields, "", rawData)
}

func dfsMask(needMaskFields map[string]MaskHandler, curPath string, rawData interface{}) interface{} {
	// 防止无限递归
	if split := strings.Split(curPath, Sep); len(split) > MaxDepth {
		return rawData
	}

	switch data := rawData.(type) {
	case map[string]interface{}:
		for key, item := range data {
			data[key] = dfsMask(needMaskFields, fmt.Sprintf("%s%s%s", curPath, Sep, key), item)
		}
	case []interface{}:
		for key, item := range data {
			data[key] = dfsMask(needMaskFields, curPath, item)
		}
	case string:
		if handler, ok := needMaskFields[curPath]; ok {
			return handler.Mask(data)
		}
	}
	return rawData
}
