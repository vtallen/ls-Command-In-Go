TARGET := vls
BUILDDIR := bin
INSTALLLOC := /usr/bin/
all:
	mkdir -p $(BUILDDIR)
	go build -o $(BUILDDIR)/$(TARGET) .
	go build -o $(TARGET)

run: 
	./$(BUILDDIR)/$(TARGET)

install: all
	sudo mv $(BUILDDIR)/$(TARGET) $(INSTALLLOC)

clean:
	rm -r $(BUILDDIR)
	touch *

