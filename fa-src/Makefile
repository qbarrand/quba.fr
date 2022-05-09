all: fa-brands.woff2 fa-solid.woff2

fa-brands.woff2: fa-brands-400.ttf unicodes-brands.txt
	pyftsubset $< --unicodes-file=unicodes-brands.txt --flavor=woff2 --output-file=$@

fa-solid.woff2: fa-solid-900.ttf unicodes-solid.txt
	pyftsubset $< --unicodes-file=unicodes-solid.txt --flavor=woff2 --output-file=$@

clean:
	rm -f fa-brands.woff2 fa-solid.woff2
