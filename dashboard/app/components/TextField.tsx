import * as Label from '@radix-ui/react-label';
import type React from 'react';

import classes from './TextField.module.css';

interface TextFieldProps {
	label: string;
	description?: string;
	error?: string;

	direction?: 'row' | 'column';
}

interface TextWrapperProps extends TextFieldProps {
	inputId: string;
	descriptionId: string;
	errorId: string;

	required?: boolean;

	children: React.ReactNode;
}

type TextInputProps = React.InputHTMLAttributes<HTMLInputElement> &
	TextFieldProps;

type TextAreaProps = React.TextareaHTMLAttributes<HTMLTextAreaElement> &
	TextFieldProps;

const TextWrapper = ({
	label,
	description,
	error,
	required,
	inputId,
	descriptionId,
	errorId,
	children,
	direction = 'column',
}: TextWrapperProps) => {
	return (
		<div className={classes.root} style={{ flexDirection: direction }}>
			<Label.Root htmlFor={inputId} className={classes.label}>
				{label}
				{required && (
					<span className={classes.required} aria-hidden="true">
						{' '}
						*
					</span>
				)}
			</Label.Root>
			{description && (
				<div id={descriptionId} className={classes.description}>
					{description}
				</div>
			)}
			{children}
			{error && (
				<div id={errorId} role="alert" className={classes.error}>
					{error}
				</div>
			)}
		</div>
	);
};

const TextInput = ({
	label,
	description,
	error,
	id,
	required,
	...props
}: TextInputProps) => {
	const inputId = id || `input-${label}`;
	const descriptionId = `${inputId}-description`;
	const errorId = `${inputId}-error`;

	return (
		<TextWrapper
			label={label}
			description={description}
			error={error}
			inputId={inputId}
			descriptionId={descriptionId}
			errorId={errorId}
			required={required}
		>
			<input
				id={inputId}
				className={classes.input}
				aria-describedby={
					`${description ? descriptionId : ''} ${error ? errorId : ''}`.trim() ||
					undefined
				}
				aria-invalid={Boolean(error)}
				data-error={error ? 'true' : undefined}
				required={required}
				{...props}
			/>
		</TextWrapper>
	);
};

const TextArea = ({
	label,
	description,
	error,
	id,
	...props
}: TextAreaProps) => {
	const textareaId = id || `textarea-${label}`;
	const descriptionId = `${textareaId}-description`;
	const errorId = `${textareaId}-error`;

	return (
		<TextWrapper
			label={label}
			description={description}
			error={error}
			inputId={textareaId}
			descriptionId={descriptionId}
			errorId={errorId}
		>
			<textarea
				id={textareaId}
				aria-describedby={
					`${description ? descriptionId : ''} ${error ? errorId : ''}`.trim() ||
					undefined
				}
				aria-invalid={!!error}
				{...props}
			/>
		</TextWrapper>
	);
};

export { TextArea, TextInput };
