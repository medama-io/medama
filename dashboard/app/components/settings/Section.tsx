import { Form } from '@remix-run/react';
import type React from 'react';

import { Button } from '@/components/Button';

import classes from './Section.module.css';

interface SectionWrapperProps {
	children: React.ReactNode;
}

interface SectionProps {
	title: string;
	description: string;
	submitDescription?: string;

	children: React.ReactNode;
	onSubmit?: () => void;
}

interface SectionDangerProps extends SectionProps {
	disabled?: boolean;
	modalChildren: React.ReactNode;
	open: () => void;
}

export const SectionWrapper = ({ children }: SectionWrapperProps) => (
	<div className={classes.sectionWrapper}>{children}</div>
);

export const SectionTitle = ({ children }: { children: React.ReactNode }) => (
	<div className={classes.title}>{children}</div>
);

export const Section = ({
	title,
	description,
	submitDescription,
	children,
	onSubmit,
}: SectionProps) => {
	return (
		<Form
			onSubmit={(e) => {
				e.preventDefault();
				e.stopPropagation();
				onSubmit?.();
			}}
		>
			<div className={classes.wrapper}>
				<SectionTitle>
					<h3>{title}</h3>
					<p style={{ marginTop: 4, marginBottom: 8 }}>{description}</p>
				</SectionTitle>
				<div className={classes.form}>{children}</div>
			</div>
			<div className={classes.divider}>
				<p>{submitDescription}</p>
				<Button className={classes.submit} type="submit">
					Save
				</Button>
			</div>
		</Form>
	);
};

// A red variant of the Section component
export const SectionDanger = ({
	title,
	description,
	submitDescription,
	children,
	modalChildren,
	disabled,
	onSubmit,
	open,
}: SectionDangerProps) => {
	return (
		<Form onSubmit={onSubmit}>
			<div className={classes.wrapper}>
				<SectionTitle>
					<h3>{title}</h3>
					<p style={{ marginTop: 4 }}>{description}</p>
				</SectionTitle>
				<div className={classes.form}>{children}</div>
			</div>
			<div className={classes.divider}>
				<p>{submitDescription}</p>
				<Button className={classes.delete} onClick={open} disabled={disabled}>
					Delete Website
				</Button>
			</div>
			{modalChildren}
		</Form>
	);
};
