package parser

import (
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Parser は、ロードしたパッケージ情報を保持する構造体です。
type Parser struct {
	pkgs map[string]*packages.Package
}

// New は、指定されたディレクトリからGoパッケージをロードし、Parserインスタンスを作成します。
// rootDirは基点となるディレクトリ、pkgDirはロード対象のパターンを指定します。
// pkgDirが空の場合は、"./..."がパターンとして使用されます。
func New(rootDir string, pkgDir string) (*Parser, error) {
	patterns := []string{"./..."}
	if pkgDir != "" {
		patterns = []string{pkgDir}
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir:   rootDir,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	parser := &Parser{
		pkgs: make(map[string]*packages.Package),
	}

	// パッケージをマップに格納
	for _, pkg := range pkgs {
		parser.pkgs[pkg.PkgPath] = pkg
	}

	return parser, nil
}

// GetAllPackages は、ロードされたすべてのパッケージを返します。
func (p *Parser) GetAllPackages() []*packages.Package {
	result := make([]*packages.Package, 0, len(p.pkgs))
	for _, pkg := range p.pkgs {
		result = append(result, pkg)
	}
	return result
}

// GetPackage は、指定されたパッケージパスのパッケージを返します。
// パッケージが見つからない場合はエラーを返します。
func (p *Parser) GetPackage(pkgPath string) (*packages.Package, error) {
	pkg, ok := p.pkgs[pkgPath]
	if !ok {
		return nil, fmt.Errorf("package not found: %s", pkgPath)
	}
	return pkg, nil
}
