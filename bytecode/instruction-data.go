package bytecode

// Data specifies the name of an instruction, and whether or not it takes
// an argument.
type Data struct {
	Name   string
	HasArg bool
}

// Instructions stores data about different instruction types.
var Instructions = map[byte]Data{
	Nop:    {Name: "NO_OP"},
	NopArg: {Name: "NO_OP_ARG", HasArg: true},

	LoadConst:      {Name: "LOAD_CONST", HasArg: true},
	LoadName:       {Name: "LOAD_NAME", HasArg: true},
	StoreName:      {Name: "STORE_NAME", HasArg: true},
	DeclareName:    {Name: "DECLARE_NAME", HasArg: true},
	LoadSubscript:  {Name: "LOAD_SUBSCRIPT"},
	StoreSubscript: {Name: "STORE_SUBSCRIPT"},

	UnaryInvert:    {Name: "UNARY_INVERT"},
	UnaryNegate:    {Name: "UNARY_NEGATE"},
	UnaryTuple:     {Name: "UNARY_TUPLE"},
	BinaryAdd:      {Name: "BINARY_ADD"},
	BinarySub:      {Name: "BINARY_SUB"},
	BinaryMul:      {Name: "BINARY_MUL"},
	BinaryDiv:      {Name: "BINARY_DIV"},
	BinaryExp:      {Name: "BINARY_EXP"},
	BinaryFloorDiv: {Name: "BINARY_FLOOR_DIV"},
	BinaryMod:      {Name: "BINARY_MODULO"},
	BinaryLogicOr:  {Name: "BINARY_LOGIC_OR"},
	BinaryLogicAnd: {Name: "BINARY_LOGIC_AND"},
	BinaryBitOr:    {Name: "BINARY_BIT_OR"},
	BinaryBitAnd:   {Name: "BINARY_BIT_AND"},
	BinaryEqual:    {Name: "BINARY_EQUAL"},
	BinaryNotEqual: {Name: "BINARY_NOT_EQUAL"},
	BinaryLess:     {Name: "BINARY_LESS_THAN"},
	BinaryMore:     {Name: "BINARY_MORE_THAN"},
	BinaryLessEq:   {Name: "BINARY_LESS_EQ"},
	BinaryMoreEq:   {Name: "BINARY_MORE_EQ"},
	BinaryTuple:    {Name: "BINARY_TUPLE"},

	CallFunction: {Name: "CALL_FUNCTION", HasArg: true},
	Return:       {Name: "RETURN"},
	PushScope:    {Name: "PUSH_SCOPE"},
	PopScope:     {Name: "POP_SCOPE"},
	Export:       {Name: "EXPORT", HasArg: true},

	Jump:       {Name: "JUMP", HasArg: true},
	JumpIf:     {Name: "JUMP_IF", HasArg: true},
	JumpUnless: {Name: "JUMP_UNLESS", HasArg: true},
	Break:      {Name: "BREAK"},
	Next:       {Name: "NEXT"},
	StartLoop:  {Name: "START_LOOP"},
	EndLoop:    {Name: "END_LOOP"},
	StartFor:   {Name: "START_FOR", HasArg: true},
	EndFor:     {Name: "END_FOR"},

	MakeList:  {Name: "MAKE_LIST", HasArg: true},
	MakeTuple: {Name: "MAKE_TUPLE", HasArg: true},
	MakeMap:   {Name: "MAKE_MAP", HasArg: true},
}
