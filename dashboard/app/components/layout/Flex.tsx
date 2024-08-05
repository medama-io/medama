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

const Group = ({ children, style }: FlexProps) => (
	<div className="group" style={style}>
		{children}
	</div>
);

export { Flex, Group };
