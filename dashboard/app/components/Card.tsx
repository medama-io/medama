import type React from 'react';
import classes from './Card.module.css';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
	children: React.ReactNode;
}

const Card = ({ children, ...props }: CardProps) => {
	return (
		<div className={classes.card} {...props}>
			{children}
		</div>
	);
};

export { Card };
