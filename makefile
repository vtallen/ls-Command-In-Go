TARGET := vls
BUILDDIR := bin

all:
	mkdir -p $(BUILDDIR)
	#go build -o $(BUILDDIR)/$(TARGET) .
	go build -o $(TARGET)

run: 
	./$(BUILDDIR)/$(TARGET)

clean:
	rm -r $(BUILDDIR)
	touch *

