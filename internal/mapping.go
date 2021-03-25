package internal

/*
// MapRowsToModel executes query and sets rows to model structure.
func MapRowsToModel(rows *sql.Rows, model interface{}) error {
	mTyp := reflect.TypeOf(model).Elem()

	// Get columns from rows.
	rCols, err := rows.Columns()
	if err != nil {
		return errors.New(err.Error(), errors.DBColumnError)
	}

	// Prepare pointers which is used to rows.Scan().
	rVal := make([][]byte, len(rCols))
	rValPtr := []interface{}{}
	for i := 0; i < len(rVal); i++ {
		rValPtr = append(rValPtr, &rVal[i])
	}

	switch mTyp.Kind() {
	case reflect.Slice, reflect.Array:
		// If model is slice of struct, generates struct whose fields are assigned
		// and appends it to vec which is reflect.Value of slice | array.
		// Or if model is slice of variable, generates variable which is assigned and
		// appends it to vec.
		vec := reflect.New(mTyp).Elem()
		if mTyp.Elem().Kind() == reflect.Struct {
			candf := ColumnsAndFields(rCols, mTyp.Elem())
			for rows.Next() {
				if err := rows.Scan(rValPtr...); err != nil {
					return errors.New(err.Error(), errors.DBScanError)
				}
				v, err := values2Fields(mTyp.Elem(), candf, rVal)
				if err != nil {
					return err
				}
				vec = reflect.Append(vec, *v)
			}
		} else {
			if len(rCols) != 1 {
				msg := fmt.Sprintf("Column length must be 1 but got %d", len(rCols))
				return errors.New(msg, errors.DBColumnError)
			}

			for rows.Next() {
				if err := rows.Scan(rValPtr...); err != nil {
					return errors.New(err.Error(), errors.DBScanError)
				}
				v, err := value2Var(mTyp.Elem(), Str(rVal[0]))
				if err != nil {
					return err
				}
				vec = reflect.Append(vec, *v)
			}
		}
		ref := reflect.ValueOf(model).Elem()
		ref.Set(vec)
		return nil
	case reflect.Struct:
		// Generates reflect.Value of struct whose fields are assigned and sets the struct to model.
		candf := ColumnsAndFields(rCols, mTyp)
		if rows.Next() {
			if err := rows.Scan(rValPtr...); err != nil {
				return errors.New(err.Error(), errors.DBScanError)
			}
			v, err := values2Fields(mTyp, candf, rVal)
			if err != nil {
				return err
			}
			ref := reflect.ValueOf(model).Elem()
			ref.Set(*v)
			return nil
		}
	case reflect.Map:
		// Generates reflect.Value of map whose key and value are assigned and sets the map to model.
		// If number of row columns is not 2, returns error.
		if len(rCols) != 2 {
			msg := fmt.Sprintf("Column length must be 2 but got %d", len(rCols))
			return errors.New(msg, errors.DBColumnError)
		}

		ref := reflect.ValueOf(model).Elem()
		for rows.Next() {
			if err := rows.Scan(rValPtr...); err != nil {
				return errors.New(err.Error(), errors.DBScanError)
			}
			k, v, err := value2Map(mTyp, Str(rVal[0]), Str(rVal[1]))
			if err != nil {
				return err
			}
			ref.SetMapIndex(*k, *v)
		}
		return nil
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool,
		reflect.String:
		// Generates reflect.Value of variable which is assigned and sets it to model.
		// If number of row columns is not 1, returns error.
		if len(rCols) != 1 {
			msg := fmt.Sprintf("Column length must be 1 but got %d", len(rCols))
			return errors.New(msg, errors.DBColumnError)
		}

		if rows.Next() {
			if err := rows.Scan(rValPtr...); err != nil {
				return errors.New(err.Error(), errors.DBScanError)
			}
			v, err := value2Var(mTyp, Str(rVal[0]))
			if err != nil {
				return err
			}
			ref := reflect.ValueOf(model).Elem()
			ref.Set(*v)
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Type %v is not supported", mTyp.Kind()), errors.InvalidTypeError)
}

// values2Fields assigns rVals to struct fields and returns the struct as pointer of reflect.Value.
// The struct is generated by mTyp which is type of model.
// vals is string values with type of [][]byte.
// candf is map of index correspondence between row columns and model fields.
func values2Fields(mTyp reflect.Type, candf map[int]int, vals [][]byte) (*reflect.Value, error) {
	v := reflect.New(mTyp).Elem()

	// Loop with row values.
	for ri := 0; ri < len(vals); ri++ {
		// fi is field index.
		fi := candf[ri]

		// Convert type of row value, []byte to string.
		val := Str(vals[ri])
		if val.Empty() {
			continue
		}

		if !v.Field(fi).CanSet() {
			msg := fmt.Sprintf("Cannot set to field %d of %s", fi, mTyp.String())
			return nil, errors.New(msg, errors.UnchangeableError)
		}

		// Field type.
		fTyp := mTyp.Field(fi).Type.Kind()

		switch fTyp {
		case reflect.String:
			s, err := val.String()
			if err != nil {
				return nil, err
			}
			v.Field(fi).SetString(s)
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			i64, err := val.Int()
			if err != nil {
				return nil, err
			}
			v.Field(fi).SetInt(i64)
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			u64, err := val.Uint()
			if err != nil {
				return nil, err
			}
			v.Field(fi).SetUint(u64)
		case reflect.Float32, reflect.Float64:
			f64, err := val.Float()
			if err != nil {
				return nil, err
			}
			v.Field(fi).SetFloat(f64)
		case reflect.Bool:
			b, err := val.Bool()
			if err != nil {
				return nil, err
			}
			v.Field(fi).SetBool(b)
		case reflect.Struct:
			if v.Field(fi).Type() == reflect.TypeOf(time.Time{}) {
				sf := mTyp.Field(fi)
				layout := TimeFormat(sf.Tag.Get("layout"))
				if layout == "" {
					layout = time.RFC3339
				}
				t, err := val.Time(layout)
				if err != nil {
					return nil, err
				}
				v.Field(fi).Set(reflect.ValueOf(t))
			}
		default:
			msg := fmt.Sprintf("%s is not supported for values2Fields", fTyp.String())
			return nil, errors.New(msg, errors.InvalidTypeError)
		}
	}

	return &v, nil
}

// value2Map assings key and val to map and returns it as pointer of reflect.Value.
// The map is generated by mTyp which is type of map.
func value2Map(mTyp reflect.Type, key Str, val Str) (*reflect.Value, *reflect.Value, error) {
	k, err := value2Var(mTyp.Key(), key)
	if err != nil {
		return nil, nil, err
	}

	v, err := value2Var(mTyp.Elem(), val)
	if err != nil {
		return nil, nil, err
	}

	return k, v, nil
}

// value2Var assigns val to variable and returns it as pointer of reflect.Value.
// The variable is generated by vTyp which is type of variable.
func value2Var(vTyp reflect.Type, val Str) (*reflect.Value, error) {
	v := reflect.New(vTyp).Elem()

	switch vTyp.Kind() {
	case reflect.String:
		s, err := val.String()
		if err != nil {
			return nil, err
		}
		v.SetString(s)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i64, err := val.Int()
		if err != nil {
			return nil, err
		}
		v.SetInt(i64)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		u64, err := val.Uint()
		if err != nil {
			return nil, err
		}
		v.SetUint(u64)
	case reflect.Float32, reflect.Float64:
		f64, err := val.Float()
		if err != nil {
			return nil, err
		}
		v.SetFloat(f64)
	case reflect.Bool:
		b, err := val.Bool()
		if err != nil {
			return nil, err
		}
		v.SetBool(b)
	case reflect.Struct:
		if vTyp == reflect.TypeOf(time.Time{}) {
			t, err := val.Time(time.RFC3339)
			if err != nil {
				return nil, err
			}
			v.Set(reflect.ValueOf(t))
		}
	default:
		msg := fmt.Sprintf("%s is not supported for values2Fields", vTyp.String())
		return nil, errors.New(msg, errors.InvalidTypeError)
	}

	return &v, nil
}
*/
