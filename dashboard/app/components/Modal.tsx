import {
	Box,
	CloseButton,
	Group,
	Text,
	TextInput,
	type TextInputProps,
	UnstyledButton,
} from '@mantine/core';
import { Form } from '@remix-run/react';
import type React from 'react';

import { IconArrowRight } from '@/components/icons/arrow-right';

import classes from './Modal.module.css';

export interface ModalProps {
	title: string;
	closeAriaLabel: string;
	description: string;
	submitLabel: string;

	onSubmit: React.FormEventHandler<HTMLFormElement>;
	resetForm: () => void;
}

export const ModalInput = (props: TextInputProps) => (
	<TextInput classNames={{ input: classes.input }} mt="md" {...props} />
);

export const ModalChild = ({
	title,
	submitLabel,
	closeAriaLabel,
	description,
	onSubmit,
	resetForm,
	children,
}: React.PropsWithChildren<ModalProps>) => {
	const resetAndClose = () => {
		resetForm();
		close();
	};

	return (
		<Box className={classes.wrapper}>
			<Group justify="space-between" align="center">
				<h2>{title}</h2>
				<CloseButton
					size="lg"
					onClick={resetAndClose}
					aria-label={closeAriaLabel}
				/>
			</Group>
			<Text size="sm" mt="xs">
				{description}
			</Text>
			<Form onSubmit={onSubmit}>
				{children}
				<UnstyledButton className={classes.submit} mt="xl" type="submit">
					<span>{submitLabel}</span>
					<IconArrowRight />
				</UnstyledButton>
			</Form>
		</Box>
	);
};
