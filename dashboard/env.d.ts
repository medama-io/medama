/// <reference types="@remix-run/node" />
/// <reference types="vite/client" />

declare module '*.mdx' {
	// biome-ignore lint/suspicious/noExplicitAny: typedef
	let MDXComponent: (props: any) => JSX.Element;
	// biome-ignore lint/suspicious/noExplicitAny: typedef
	export const frontmatter: any;
	export default MDXComponent;
}
