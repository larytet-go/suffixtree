# Generalized Suffix Tree
A Go implementation of a Generalized Suffix Tree using Ukkonen's algorithm

The package just translate from Alessandro Bahgat Shehata's java version to golang and do some optimization  
For more details, you should look at [abahgat/suffixtree](https://github.com/abahgat/suffixtree/) 

## Usage

```go
package main

import (
	"fmt"

	"github.com/ljfuyuan/suffixtree"
)

func main() {
	words := []string{"banana", "apple", "中文app"}
	tree := suffixtree.NewGeneralizedSuffixTree()
	for k, word := range words {
		tree.Put(word, k)
	}
	indexs := tree.Search("a", -1)

	fmt.Println(indexs)
	//[0 2 1]
	for _, index := range indexs {
		fmt.Println(words[index])
	}
	//banana
	//中文app
	//apple
}
```
## Reading 


* https://stackoverflow.com/questions/14032903/longest-common-contiguous-subsequence-algorithm
* https://en.wikipedia.org/wiki/Suffix_tree
* https://en.wikipedia.org/wiki/Suffix_array
* https://github.com/ljfuyuan/suffixtree
* https://github.com/cweill/SuffixArray-Golang
* https://www.sciencedirect.com/science/article/pii/S1570866709000343
* https://en.wikipedia.org/wiki/Suffix_array


## License

This Generalized Suffix Tree is released under the Apache License 2.0

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
