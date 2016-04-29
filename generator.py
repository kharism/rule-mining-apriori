import random

rand = random.Random()
for i in range(40):
	s = str(i+1)+", "
	item = []
	for j in range(random.randint(3,20)):
		u = random.randint(0,9)
		if not u in item:
			item.append(u)
	for u in item:
		s+=str(u)+", "
	print(s)
