// Code generated by ogen, DO NOT EDIT.

package api

// setDefaults set default value of fields.
func (s *BadRequestErrorError) setDefaults() {
	{
		val := int32(400)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *ConflictErrorError) setDefaults() {
	{
		val := int32(409)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *ForbiddenErrorError) setDefaults() {
	{
		val := int32(403)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *InternalServerErrorError) setDefaults() {
	{
		val := int32(500)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *NotFoundErrorError) setDefaults() {
	{
		val := int32(404)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *UnauthorisedErrorError) setDefaults() {
	{
		val := int32(401)
		s.Code = val
	}
}

// setDefaults set default value of fields.
func (s *UserSettings) setDefaults() {
	{
		val := UserSettingsLanguage("en")
		s.Language.SetTo(val)
	}
	{
		val := UserSettingsScriptType("default")
		s.ScriptType.SetTo(val)
	}
}
