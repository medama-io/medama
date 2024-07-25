import type React from 'react';
import type { CSSProperties } from 'react';

import clsx from 'clsx';
import classes from './Flex.module.css';

interface FlexProps {
	children: React.ReactNode;
	className?: string;

	align?: CSSProperties['alignItems'];
	justify?: CSSProperties['justifyContent'];
	gap?: CSSProperties['gap'];
}

const Flex = ({ children, className, justify, align, gap }: FlexProps) => (
	<div
		className={clsx(classes.flex, className)}
		style={{ justifyContent: justify, alignItems: align, gap }}
	>
		{children}
	</div>
);

const Group = ({ children, className, justify, align, gap }: FlexProps) => (
	<div
		className={clsx(classes.group, className)}
		style={{ justifyContent: justify, alignItems: align, gap }}
	>
		{children}
	</div>
);

export { Flex, Group };
