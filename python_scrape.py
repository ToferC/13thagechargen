import workflow
from bs4 import BeautifulSoup
import requests

class_list = ['barbarian', 'bard', 'cleric', 'fighter', 'paladin', 'ranger', 'rogue', 'sorcerer', 'wizard']

class_select = workflow.get_variable('class_select')

item_url = "http://www.13thagesrd.com/classes/{}".format(class_select.lower())

r = requests.get(item_url)

data = r.text

soup = BeautifulSoup(data)

source = {}

ability_soup = soup.find_all('h4')

void_abilities = ['', 'Melee Weapons', 'Ranged Weapons', 'Armor', 'Gold Pieces']

for section in ability_soup:
	if section.text in void_abilities:
		pass
	else:
	#print(section.text)
		source[section.text] = []
		nextNode = section
		while True:
			nextNode = nextNode.find_next_sibling()
			try:
				tag_name = nextNode.name
			except AttributeError:
				tag_name = ""
			if tag_name in ["h5", "p", 't']:
				source[section.text].append(nextNode.text)
				#print nextNode.text
			else:
				#print "*****"
				break

source_txt = "{}".format(source)

source_list = ''
for k in source.keys():
	source_list += "{}\n".format(k)

workflow.set_output(source_list)
workflow.set_variable('source',source_txt)
