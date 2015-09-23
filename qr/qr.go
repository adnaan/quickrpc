package qr

import (
	"path"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

const (
	qrSQLPkgPath = "github.com/adnaan/quickrpc/sql"
)

func init() {
	generator.RegisterPlugin(new(qr))
}

// qr is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for qr support.
type qr struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "qr".
func (q *qr) Name() string {
	return "qr"
}

var contextPkg = "context"
var qrSQLPkg string

// Init initializes the plugin.
func (q *qr) Init(gen *generator.Generator) {
	q.gen = gen
	qrSQLPkg = generator.RegisterUniquePackageName("sql", nil)
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (q *qr) objectNamed(name string) generator.Object {
	q.gen.RecordTypeUse(name)
	return q.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (q *qr) typeName(str string) string {
	return q.gen.TypeName(q.objectNamed(str))
}

// P forwards to q.gen.P.
func (q *qr) P(args ...interface{}) { q.gen.P(args...) }

// Generate generates code for the services in the given file.
func (q *qr) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	// q.P("// Reference imports to suppress errors if they are not otherwise used.")
	// q.P("var _ ", contextPkg, ".Context")
	// q.P("var _ ", qrPkg, ".ClientConn")
	// q.P()
	for i, service := range file.FileDescriptorProto.Service {
		q.generateServiceImpl(file, service, i)
	}

	for _, message := range file.FileDescriptorProto.MessageType {
		name := message.GetName()

		options := message.GetOptions()
		if options != nil {
			q.P()
			q.P("// Generate sql queries: ", name, " options ", options.String())
			q.P()

		}

	}
}

// GenerateImports generates the import declaration for this file.
func (q *qr) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	q.P("import (")
	q.P(qrSQLPkg, " ", strconv.Quote(path.Join(q.gen.ImportPrefix, qrSQLPkgPath)))
	q.P(")")
	q.P()
}

// generateServiceImpl generates implementation for the named service.
func (q *qr) generateServiceImpl(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {

	serverType := "Server"
	origServName := service.GetName()
	servName := generator.CamelCase(origServName)
	q.P()
	q.P("// Generated Implementation for ", servName, " service")
	q.P()

	q.P("type ", serverType, " struct {")
	// TODO server fields
	q.P("}")
	q.P()

	for _, method := range service.Method {
		q.generateServiceMethod(serverType, servName, method)
	}

}

func (q *qr) generateServiceMethod(serverType, servName string, method *pb.MethodDescriptorProto) {
	q.P("func (s *", serverType, ") ", q.generateMethodSignature(servName, method), "{")
	//TODO method body
	q.P("return &", q.typeName(method.GetOutputType()), "{}, nil")
	q.P("}")
	q.P()
}

// generateServerSignature returns the server-side signature for a method.
func (q *qr) generateMethodSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)

	var reqArgs []string
	ret := "error"
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		reqArgs = append(reqArgs, contextPkg+".Context")
		ret = "(*" + q.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "*"+q.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Server")
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}
