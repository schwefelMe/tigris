// Copyright 2022-2023 Tigris Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint:dupl
package schema

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestJavaSchemaGenerator(t *testing.T) {
	cases := []struct {
		name string
		in   string
		exp  string
	}{
		{
			"types", typesTest, `
class Product {
    private long[] arrInts;
    private boolean bool;
    private byte[] byte1;
    private int id;
    private long int64;
    @TigrisField(description = "field description")
    private long int64WithDesc;
    private String name;
    private double price;
    private Date time1;
    private long[] twoDArr;
    private UUID uUID1;

    public long[] getArrInts() {
        return arrInts;
    }

    public void setArrInts(long[] arrInts) {
        this.arrInts = arrInts;
    }

    public boolean isBool() {
        return bool;
    }

    public void setBool(boolean bool) {
        this.bool = bool;
    }

    public byte[] getByte1() {
        return byte1;
    }

    public void setByte1(byte[] byte1) {
        this.byte1 = byte1;
    }

    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public long getInt64() {
        return int64;
    }

    public void setInt64(long int64) {
        this.int64 = int64;
    }

    public long getInt64WithDesc() {
        return int64WithDesc;
    }

    public void setInt64WithDesc(long int64WithDesc) {
        this.int64WithDesc = int64WithDesc;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public double getPrice() {
        return price;
    }

    public void setPrice(double price) {
        this.price = price;
    }

    public Date getTime1() {
        return time1;
    }

    public void setTime1(Date time1) {
        this.time1 = time1;
    }

    public long[] getTwoDArr() {
        return twoDArr;
    }

    public void setTwoDArr(long[] twoDArrs) {
        this.twoDArr = twoDArrs;
    }

    public UUID getUUID1() {
        return uUID1;
    }

    public void setUUID1(UUID uUID1) {
        this.uUID1 = uUID1;
    }

    public Product() {};

    public Product(
        long[] arrInts,
        boolean bool,
        byte[] byte1,
        int id,
        long int64,
        long int64WithDesc,
        String name,
        double price,
        Date time1,
        long[] twoDArrs,
        UUID uUID1
    ) {
        this.arrInts = arrInts;
        this.bool = bool;
        this.byte1 = byte1;
        this.id = id;
        this.int64 = int64;
        this.int64WithDesc = int64WithDesc;
        this.name = name;
        this.price = price;
        this.time1 = time1;
        this.twoDArr = twoDArrs;
        this.uUID1 = uUID1;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        Product other = (Product) o;
        return
            Arrays.equals(arrInts, other.arrInts) &&
            bool == other.bool &&
            byte1 == other.byte1 &&
            id == other.id &&
            int64 == other.int64 &&
            int64WithDesc == other.int64WithDesc &&
            name == other.name &&
            price == other.price &&
            time1 == other.time1 &&
            Arrays.equals(twoDArr, other.twoDArr) &&
            uUID1 == other.uUID1;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            arrInts,
            bool,
            byte1,
            id,
            int64,
            int64WithDesc,
            name,
            price,
            time1,
            twoDArr,
            uUID1
        );
    }
}
`,
		},
		{
			"tags", tagsTest, `
// Product type description
@com.tigrisdata.db.annotation.TigrisCollection(value = "products")
public class Product implements TigrisCollectionType {
    @TigrisPrimaryKey(autoGenerate = true)
    private int Gen;
    @TigrisPrimaryKey(order = 1)
    private int Key;
    @TigrisPrimaryKey(order = 2, autoGenerate = true)
    private int KeyGenIdx;
    private String def_val_cuid;
    private Date def_val_date;
    private Date def_val_date_const;
    private long def_val_int;
    private String def_val_str;
    private String def_val_str_q;
    private UUID def_val_uuid;
    private String max_len_str;
    private String max_len_str_req;
    @TigrisPrimaryKey(autoGenerate = true)
    private int name_gen;
    @TigrisPrimaryKey(order = 4, autoGenerate = true)
    private int name_gen_key;
    @TigrisPrimaryKey(order = 3)
    private int name_key;
    private String name_si;
    private int name_sif;
    private Date name_sifs;
    private int req_field;
    private Date time_f;
    private int user_name;

    public int getGen() {
        return Gen;
    }

    public void setGen(int gen) {
        this.Gen = gen;
    }

    public int getKey() {
        return Key;
    }

    public void setKey(int key) {
        this.Key = key;
    }

    public int getKeyGenIdx() {
        return KeyGenIdx;
    }

    public void setKeyGenIdx(int keyGenIdx) {
        this.KeyGenIdx = keyGenIdx;
    }

    public String getDef_val_cuid() {
        return def_val_cuid;
    }

    public void setDef_val_cuid(String defValCuid) {
        this.def_val_cuid = defValCuid;
    }

    public Date getDef_val_date() {
        return def_val_date;
    }

    public void setDef_val_date(Date defValDate) {
        this.def_val_date = defValDate;
    }

    public Date getDef_val_date_const() {
        return def_val_date_const;
    }

    public void setDef_val_date_const(Date defValDateConst) {
        this.def_val_date_const = defValDateConst;
    }

    public long getDef_val_int() {
        return def_val_int;
    }

    public void setDef_val_int(long defValInt) {
        this.def_val_int = defValInt;
    }

    public String getDef_val_str() {
        return def_val_str;
    }

    public void setDef_val_str(String defValStr) {
        this.def_val_str = defValStr;
    }

    public String getDef_val_str_q() {
        return def_val_str_q;
    }

    public void setDef_val_str_q(String defValStrQ) {
        this.def_val_str_q = defValStrQ;
    }

    public UUID getDef_val_uuid() {
        return def_val_uuid;
    }

    public void setDef_val_uuid(UUID defValUuid) {
        this.def_val_uuid = defValUuid;
    }

    public String getMax_len_str() {
        return max_len_str;
    }

    public void setMax_len_str(String maxLenStr) {
        this.max_len_str = maxLenStr;
    }

    public String getMax_len_str_req() {
        return max_len_str_req;
    }

    public void setMax_len_str_req(String maxLenStrReq) {
        this.max_len_str_req = maxLenStrReq;
    }

    public int getName_gen() {
        return name_gen;
    }

    public void setName_gen(int nameGen) {
        this.name_gen = nameGen;
    }

    public int getName_gen_key() {
        return name_gen_key;
    }

    public void setName_gen_key(int nameGenKey) {
        this.name_gen_key = nameGenKey;
    }

    public int getName_key() {
        return name_key;
    }

    public void setName_key(int nameKey) {
        this.name_key = nameKey;
    }

    public String getName_si() {
        return name_si;
    }

    public void setName_si(String nameSi) {
        this.name_si = nameSi;
    }

    public int getName_sif() {
        return name_sif;
    }

    public void setName_sif(int nameSif) {
        this.name_sif = nameSif;
    }

    public Date getName_sifs() {
        return name_sifs;
    }

    public void setName_sifs(Date nameSifs) {
        this.name_sifs = nameSifs;
    }

    public int getReq_field() {
        return req_field;
    }

    public void setReq_field(int reqField) {
        this.req_field = reqField;
    }

    public Date getTime_f() {
        return time_f;
    }

    public void setTime_f(Date timeF) {
        this.time_f = timeF;
    }

    public int getUser_name() {
        return user_name;
    }

    public void setUser_name(int userName) {
        this.user_name = userName;
    }

    public Product() {};

    public Product(
        int gen,
        int key,
        int keyGenIdx,
        String defValCuid,
        Date defValDate,
        Date defValDateConst,
        long defValInt,
        String defValStr,
        String defValStrQ,
        UUID defValUuid,
        String maxLenStr,
        String maxLenStrReq,
        int nameGen,
        int nameGenKey,
        int nameKey,
        String nameSi,
        int nameSif,
        Date nameSifs,
        int reqField,
        Date timeF,
        int userName
    ) {
        this.Gen = gen;
        this.Key = key;
        this.KeyGenIdx = keyGenIdx;
        this.def_val_cuid = defValCuid;
        this.def_val_date = defValDate;
        this.def_val_date_const = defValDateConst;
        this.def_val_int = defValInt;
        this.def_val_str = defValStr;
        this.def_val_str_q = defValStrQ;
        this.def_val_uuid = defValUuid;
        this.max_len_str = maxLenStr;
        this.max_len_str_req = maxLenStrReq;
        this.name_gen = nameGen;
        this.name_gen_key = nameGenKey;
        this.name_key = nameKey;
        this.name_si = nameSi;
        this.name_sif = nameSif;
        this.name_sifs = nameSifs;
        this.req_field = reqField;
        this.time_f = timeF;
        this.user_name = userName;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        Product other = (Product) o;
        return
            Gen == other.Gen &&
            Key == other.Key &&
            KeyGenIdx == other.KeyGenIdx &&
            def_val_cuid == other.def_val_cuid &&
            def_val_date == other.def_val_date &&
            def_val_date_const == other.def_val_date_const &&
            def_val_int == other.def_val_int &&
            def_val_str == other.def_val_str &&
            def_val_str_q == other.def_val_str_q &&
            def_val_uuid == other.def_val_uuid &&
            max_len_str == other.max_len_str &&
            max_len_str_req == other.max_len_str_req &&
            name_gen == other.name_gen &&
            name_gen_key == other.name_gen_key &&
            name_key == other.name_key &&
            name_si == other.name_si &&
            name_sif == other.name_sif &&
            name_sifs == other.name_sifs &&
            req_field == other.req_field &&
            time_f == other.time_f &&
            user_name == other.user_name;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            Gen,
            Key,
            KeyGenIdx,
            def_val_cuid,
            def_val_date,
            def_val_date_const,
            def_val_int,
            def_val_str,
            def_val_str_q,
            def_val_uuid,
            max_len_str,
            max_len_str_req,
            name_gen,
            name_gen_key,
            name_key,
            name_si,
            name_sif,
            name_sifs,
            req_field,
            time_f,
            user_name
        );
    }
}
`,
		},
		{"object", objectTest, `
class SubArrayNested {
    private int field_arr;

    public int getField_arr() {
        return field_arr;
    }

    public void setField_arr(int fieldArr) {
        this.field_arr = fieldArr;
    }

    public SubArrayNested() {};

    public SubArrayNested(
        int fieldArr
    ) {
        this.field_arr = fieldArr;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubArrayNested other = (SubArrayNested) o;
        return
            field_arr == other.field_arr;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_arr
        );
    }
}

class SubObjectNested {
    private int field_3;

    public int getField_3() {
        return field_3;
    }

    public void setField_3(int field3) {
        this.field_3 = field3;
    }

    public SubObjectNested() {};

    public SubObjectNested(
        int field3
    ) {
        this.field_3 = field3;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubObjectNested other = (SubObjectNested) o;
        return
            field_3 == other.field_3;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_3
        );
    }
}

class SubObjectNestedOne {
    private int field_31;
    private SubObjectNested subObjectNested;

    public int getField_31() {
        return field_31;
    }

    public void setField_31(int field31) {
        this.field_31 = field31;
    }

    public SubObjectNested getSubObjectNested() {
        return subObjectNested;
    }

    public void setSubObjectNested(SubObjectNested subObjectNested) {
        this.subObjectNested = subObjectNested;
    }

    public SubObjectNestedOne() {};

    public SubObjectNestedOne(
        int field31,
        SubObjectNested subObjectNested
    ) {
        this.field_31 = field31;
        this.subObjectNested = subObjectNested;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubObjectNestedOne other = (SubObjectNestedOne) o;
        return
            field_31 == other.field_31 &&
            Objects.equals(subObjectNested, other.subObjectNested);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_31,
            subObjectNested
        );
    }
}

class SubObjectNested1 {
    private int field_4;

    public int getField_4() {
        return field_4;
    }

    public void setField_4(int field4) {
        this.field_4 = field4;
    }

    public SubObjectNested1() {};

    public SubObjectNested1(
        int field4
    ) {
        this.field_4 = field4;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubObjectNested1 other = (SubObjectNested1) o;
        return
            field_4 == other.field_4;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_4
        );
    }
}

class SubObjectNestedThree {
    private int field_33;
    private SubObjectNested1 subObjectNested;

    public int getField_33() {
        return field_33;
    }

    public void setField_33(int field33) {
        this.field_33 = field33;
    }

    public SubObjectNested1 getSubObjectNested() {
        return subObjectNested;
    }

    public void setSubObjectNested(SubObjectNested1 subObjectNested) {
        this.subObjectNested = subObjectNested;
    }

    public SubObjectNestedThree() {};

    public SubObjectNestedThree(
        int field33,
        SubObjectNested1 subObjectNested
    ) {
        this.field_33 = field33;
        this.subObjectNested = subObjectNested;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubObjectNestedThree other = (SubObjectNestedThree) o;
        return
            field_33 == other.field_33 &&
            Objects.equals(subObjectNested, other.subObjectNested);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_33,
            subObjectNested
        );
    }
}

class SubObjectNestedTwo {
    private int field_32;
    private SubObjectNested1 subObjectNested;

    public int getField_32() {
        return field_32;
    }

    public void setField_32(int field32) {
        this.field_32 = field32;
    }

    public SubObjectNested1 getSubObjectNested() {
        return subObjectNested;
    }

    public void setSubObjectNested(SubObjectNested1 subObjectNested) {
        this.subObjectNested = subObjectNested;
    }

    public SubObjectNestedTwo() {};

    public SubObjectNestedTwo(
        int field32,
        SubObjectNested1 subObjectNested
    ) {
        this.field_32 = field32;
        this.subObjectNested = subObjectNested;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubObjectNestedTwo other = (SubObjectNestedTwo) o;
        return
            field_32 == other.field_32 &&
            Objects.equals(subObjectNested, other.subObjectNested);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_32,
            subObjectNested
        );
    }
}

class SubArray {
    private int field_3;
    private SubArrayNested[] subArrayNesteds;
    private SubObjectNested subObjectNested;
    private SubObjectNestedOne subObjectNestedOne;
    private SubObjectNested subObjectNestedReuseTypeByBody;
    private SubObjectNestedThree subObjectNestedThree;
    private SubObjectNestedTwo subObjectNestedTwo;

    public int getField_3() {
        return field_3;
    }

    public void setField_3(int field3) {
        this.field_3 = field3;
    }

    public SubArrayNested[] getSubArrayNesteds() {
        return subArrayNesteds;
    }

    public void setSubArrayNesteds(SubArrayNested[] subArrayNesteds) {
        this.subArrayNesteds = subArrayNesteds;
    }

    public SubObjectNested getSubObjectNested() {
        return subObjectNested;
    }

    public void setSubObjectNested(SubObjectNested subObjectNested) {
        this.subObjectNested = subObjectNested;
    }

    public SubObjectNestedOne getSubObjectNestedOne() {
        return subObjectNestedOne;
    }

    public void setSubObjectNestedOne(SubObjectNestedOne subObjectNestedOne) {
        this.subObjectNestedOne = subObjectNestedOne;
    }

    public SubObjectNested getSubObjectNestedReuseTypeByBody() {
        return subObjectNestedReuseTypeByBody;
    }

    public void setSubObjectNestedReuseTypeByBody(SubObjectNested subObjectNestedReuseTypeByBody) {
        this.subObjectNestedReuseTypeByBody = subObjectNestedReuseTypeByBody;
    }

    public SubObjectNestedThree getSubObjectNestedThree() {
        return subObjectNestedThree;
    }

    public void setSubObjectNestedThree(SubObjectNestedThree subObjectNestedThree) {
        this.subObjectNestedThree = subObjectNestedThree;
    }

    public SubObjectNestedTwo getSubObjectNestedTwo() {
        return subObjectNestedTwo;
    }

    public void setSubObjectNestedTwo(SubObjectNestedTwo subObjectNestedTwo) {
        this.subObjectNestedTwo = subObjectNestedTwo;
    }

    public SubArray() {};

    public SubArray(
        int field3,
        SubArrayNested[] subArrayNesteds,
        SubObjectNested subObjectNested,
        SubObjectNestedOne subObjectNestedOne,
        SubObjectNested subObjectNestedReuseTypeByBody,
        SubObjectNestedThree subObjectNestedThree,
        SubObjectNestedTwo subObjectNestedTwo
    ) {
        this.field_3 = field3;
        this.subArrayNesteds = subArrayNesteds;
        this.subObjectNested = subObjectNested;
        this.subObjectNestedOne = subObjectNestedOne;
        this.subObjectNestedReuseTypeByBody = subObjectNestedReuseTypeByBody;
        this.subObjectNestedThree = subObjectNestedThree;
        this.subObjectNestedTwo = subObjectNestedTwo;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        SubArray other = (SubArray) o;
        return
            field_3 == other.field_3 &&
            Arrays.equals(subArrayNesteds, other.subArrayNesteds) &&
            Objects.equals(subObjectNested, other.subObjectNested) &&
            Objects.equals(subObjectNestedOne, other.subObjectNestedOne) &&
            Objects.equals(subObjectNestedReuseTypeByBody, other.subObjectNestedReuseTypeByBody) &&
            Objects.equals(subObjectNestedThree, other.subObjectNestedThree) &&
            Objects.equals(subObjectNestedTwo, other.subObjectNestedTwo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            field_3,
            subArrayNesteds,
            subObjectNested,
            subObjectNestedOne,
            subObjectNestedReuseTypeByBody,
            subObjectNestedThree,
            subObjectNestedTwo
        );
    }
}

// Subtype sub type description
class Subtype {
    private int id2;

    public int getId2() {
        return id2;
    }

    public void setId2(int id2) {
        this.id2 = id2;
    }

    public Subtype() {};

    public Subtype(
        int id2
    ) {
        this.id2 = id2;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        Subtype other = (Subtype) o;
        return
            id2 == other.id2;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            id2
        );
    }
}

@com.tigrisdata.db.annotation.TigrisCollection(value = "products")
public class Product implements TigrisCollectionType {
    private SubArray[] subArrays;
    @TigrisField(description = "sub type description")
    private Subtype subtype;

    public SubArray[] getSubArrays() {
        return subArrays;
    }

    public void setSubArrays(SubArray[] subArrays) {
        this.subArrays = subArrays;
    }

    public Subtype getSubtype() {
        return subtype;
    }

    public void setSubtype(Subtype subtype) {
        this.subtype = subtype;
    }

    public Product() {};

    public Product(
        SubArray[] subArrays,
        Subtype subtype
    ) {
        this.subArrays = subArrays;
        this.subtype = subtype;
    };

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        Product other = (Product) o;
        return
            Arrays.equals(subArrays, other.subArrays) &&
            Objects.equals(subtype, other.subtype);
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            subArrays,
            subtype
        );
    }
}
`},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			w := bufio.NewWriter(&buf)

			s := schemaGenerator{
				hasTime:     false,
				hasUUID:     false,
				langTypeGen: &JSONToJava{},
				writer:      w,
				types:       make(map[string][]string),
				bodyToType:  make(map[string]string),
			}

			err := s.genCollectionSchema([]byte(v.in))
			require.NoError(t, err)

			_ = w.Flush()
			assert.Equal(t, v.exp, buf.String())
		})
	}
}
