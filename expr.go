package eep

type Visitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitLogicalExpr(expr *Logical) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitCallExpr(expr *Call) interface{}
	VisitVariableExpr(expr *Variable) interface{}
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func NewBinary(left, right Expr, operator *Token) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func (bin *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(bin)
}

type Call struct {
	Callee    Expr
	Paren     *Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *Token, arguments []Expr) *Call {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (cl *Call) Accept(v Visitor) interface{} {
	return v.VisitCallExpr(cl)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{Expression: expr}
}

func (g *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(v interface{}) *Literal {
	return &Literal{Value: v}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(l)
}

type Logical struct {
	Left     Expr
	Right    Expr
	Operator *Token
}

func NewLogical(left, right Expr, op *Token) *Logical {
	return &Logical{Left: left, Right: right, Operator: op}
}

func (l *Logical) Accept(v Visitor) interface{} {
	return v.VisitLogicalExpr(l)
}

type Unary struct {
	Operator *Token
	Right    Expr
}

func NewUnary(operator *Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(u)
}

type Variable struct {
	Name *Token
}

func NewVariable(name *Token) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Accept(vt Visitor) interface{} {
	return vt.VisitVariableExpr(v)
}
