package generator

// credit to github.com/vektra/mockery

import (
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
)

var invalidIdentifierChar = regexp.MustCompile("[^[:digit:][:alpha:]_]")

// copied from github.com/vektra/mockery
func (vo *VanillaMockStructOutput) renderTypeTuple(tup *types.Tuple) string {
	var parts []string

	for i := 0; i < tup.Len(); i++ {
		v := tup.At(i)
		parts = append(parts, vo.renderType(v.Type()))
	}

	return strings.Join(parts, " , ")
}

// copied from github.com/vektra/mockery
func (v *VanillaMockStructOutput) renderType(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Named:
		o := t.Obj()
		if o.Pkg() == nil || o.Pkg().Name() == "main" { // || (!KeepTree && InPackage && o.Pkg() == iface.Pkg) {
			return o.Name()
		}
		return v.addPackageImport(o.Pkg()) + "." + o.Name()
	case *types.Basic:
		return t.Name()
	case *types.Pointer:
		return "*" + v.renderType(t.Elem())
	case *types.Slice:
		return "[]" + v.renderType(t.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), v.renderType(t.Elem()))
	case *types.Signature:
		switch t.Results().Len() {
		case 0:
			return fmt.Sprintf(
				"func(%s)",
				v.renderTypeTuple(t.Params()),
			)
		case 1:
			return fmt.Sprintf(
				"func(%s) %s",
				v.renderTypeTuple(t.Params()),
				v.renderType(t.Results().At(0).Type()),
			)
		default:
			return fmt.Sprintf(
				"func(%s)(%s)",
				v.renderTypeTuple(t.Params()),
				v.renderTypeTuple(t.Results()),
			)
		}
	case *types.Map:
		kt := v.renderType(t.Key())
		vt := v.renderType(t.Elem())

		return fmt.Sprintf("map[%s]%s", kt, vt)
	case *types.Chan:
		switch t.Dir() {
		case types.SendRecv:
			return "chan " + v.renderType(t.Elem())
		case types.RecvOnly:
			return "<-chan " + v.renderType(t.Elem())
		default:
			return "chan<- " + v.renderType(t.Elem())
		}
	case *types.Struct:
		var fields []string

		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)

			if f.Anonymous() {
				fields = append(fields, v.renderType(f.Type()))
			} else {
				fields = append(fields, fmt.Sprintf("%s %s", f.Name(), v.renderType(f.Type())))
			}
		}

		return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
	case *types.Interface:
		if t.NumMethods() != 0 {
			panic("Unable to mock inline interfaces with methods")
		}

		return "interface{}"
	default:
		panic(fmt.Sprintf("un-namable type: %#v (%T)", t, t))
	}
}

func (v *VanillaMockStructOutput) getNonConflictingName(path string, name string) string {
	if !v.importNameExists(name) {
		return name
	}

	// The path will always contain '/' because it is enforced in getLocalizedPath
	// regardless of OS.
	directories := strings.Split(path, "/")

	cleanedDirectories := make([]string, 0, len(directories))
	for _, directory := range directories {
		cleaned := invalidIdentifierChar.ReplaceAllString(directory, "_")
		cleanedDirectories = append(cleanedDirectories, cleaned)
	}
	numDirectories := len(cleanedDirectories)
	var prospectiveName string
	for i := 1; i <= numDirectories; i++ {
		prospectiveName = strings.Join(cleanedDirectories[numDirectories-i:], "")
		if !v.importNameExists(prospectiveName) {
			return prospectiveName
		}
	}
	// Try adding numbers to the given name
	i := 2
	for {
		prospectiveName = fmt.Sprintf("%v%d", name, i)
		if !v.importNameExists(prospectiveName) {
			return prospectiveName
		}
		i++
	}
}
func (v *VanillaMockStructOutput) addPackageImportWithName(path, name string) string {
	path = v.getLocalizedPath(path)
	if existingName, pathExists := v.packagePathToName[path]; pathExists {
		return existingName
	}

	nonConflictingName := v.getNonConflictingName(path, name)
	v.packagePathToName[path] = nonConflictingName
	v.nameToPackagePath[nonConflictingName] = path
	return nonConflictingName
}

func (v *VanillaMockStructOutput) getLocalizedPath(path string) string {

	if strings.HasSuffix(path, ".go") {
		path, _ = filepath.Split(path)
	}
	if localized, ok := v.localizationCache[path]; ok {
		return localized
	}
	directories := strings.Split(path, string(filepath.Separator))
	numDirectories := len(directories)
	vendorIndex := -1
	for i := 1; i <= numDirectories; i++ {
		dir := directories[numDirectories-i]
		if dir == "vendor" {
			vendorIndex = numDirectories - i
			break
		}
	}

	toReturn := path
	if vendorIndex >= 0 {
		toReturn = filepath.Join(directories[vendorIndex+1:]...)
	} else if filepath.IsAbs(path) {
		toReturn = calculateImport(v.packageRoots, path)
	}

	// Enforce '/' slashes for import paths in every OS.
	toReturn = filepath.ToSlash(toReturn)

	v.localizationCache[path] = toReturn
	return toReturn
}

func calculateImport(set []string, path string) string {
	for _, root := range set {
		if strings.HasPrefix(path, root) {
			packagePath, err := filepath.Rel(root, path)
			if err == nil {
				return packagePath
			}
		}
	}
	return path
}
