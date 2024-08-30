import type React from 'react';
import type { CSSProperties } from 'react';

interface FlexProps {
	children: React.ReactNode;
	style?: CSSProperties;
}

const Flex = ({ children, style }: FlexProps) => (
	<div className="flex" style={style}>
		{children}
	</div>
);

const Group = ({ children, style, ...rest }: FlexProps) => (
	<div className="group" style={style} {...rest}>
		{children}
	</div>
);

export { Flex, Group };
