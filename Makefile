
.PHONEY:deploy-local
deploy-local: package
	adb install dist/bowlsscore.apk


.PHONEY: package
package:
	mkdir -p dist
	fyne package -os android
	mv Bowls_Scorer.apk dist/bowlsscore.apk

.PHONEY: build
build:
	mkdir -p bin
	fyne build -o bin/bowlscore
.PHONEY: run
run: build
	bin/bowlscore