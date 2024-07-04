import { Stack, Text, Title, UnstyledButton } from '@mantine/core';
import { Form } from '@remix-run/react';
import clsx from 'clsx';
import type React from 'react';

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
		<Form onSubmit={onSubmit}>
			<div className={classes.wrapper}>
				<SectionTitle>
					<Title order={3}>{title}</Title>
					<Text mt="xs">{description}</Text>
				</SectionTitle>
				<Stack className={classes.form}>{children}</Stack>
			</div>
			<div className={classes.divider}>
				<Text>{submitDescription}</Text>
				<UnstyledButton className={classes.submit} type="submit">
					Save
				</UnstyledButton>
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
	onSubmit,
	open,
}: SectionDangerProps) => {
	return (
		<Form onSubmit={onSubmit}>
			<div className={classes.wrapper}>
				<SectionTitle>
					<Title order={3}>{title}</Title>
					<Text mt="xs">{description}</Text>
				</SectionTitle>
				<Stack className={classes.form}>{children}</Stack>
			</div>
			<div className={classes.divider}>
				<Text>{submitDescription}</Text>
				<UnstyledButton
					className={clsx(classes.submit, classes.delete)}
					onClick={open}
				>
					Delete Website
				</UnstyledButton>
			</div>
			{modalChildren}
		</Form>
	);
};
