package main

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

var examples = []string{
	`<para>TEXT1 <ref refid="ID" kindref="member">REF</ref> TEXT2</para>`,
	`<para><bold>Example</bold> <programlisting><codeline><highlight class="normal">typedef<sp/>struct<sp/>mcs_node_s<sp/>{</highlight></codeline>
	 <codeline><highlight class="normal"><sp/><sp/><sp/><sp/>vatomicptr(struct<sp/>mcs_node_s*)<sp/>next;</highlight></codeline>
	<codeline><highlight class="normal">}<sp/>mcs_node_t;</highlight></codeline>
	 </programlisting> </para>`,
}

func Test_ParaTest(t *testing.T) {
	for i, r := range examples {
		fmt.Println("TEST ", i)
		reader := strings.NewReader(r)
		var para Para
		err := xml.NewDecoder(reader).Decode(&para)
		if err != nil {
			fmt.Println(err)
			t.Errorf("could not decode xml: %v", err)
			t.FailNow()
		}
		fmt.Println()
		//spew.Dump("Parsing:", para)
	}
}

// func Test_ParaDump(t *testing.T) {
// 	for i, r := range examples {
// 		fmt.Println("TEST ", i, r)
// 	}
// }
