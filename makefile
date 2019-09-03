#  Copyright 2019 Marius Ackerman
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

.PHONY: clean lexer gocc debug install

install: clean gocc
	go install

clean:
	rm first.txt ; \
	rm lexer_sets.txt ; \
	rm LR1_* ; \
	rm terminals.txt ; \
	rm -rf lexer ; \
	rm -rf parser ; \

gocc:
	gocc -p gogll gogll.bnf 

debug:
	gocc -p gogll -v -debug_lexer gogll.bnf

