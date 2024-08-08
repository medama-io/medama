import { EyeNoneIcon, EyeOpenIcon } from '@radix-ui/react-icons';
import * as Label from '@radix-ui/react-label';
import type React from 'react';

import { IconButton } from '@/components/Button';
import { useDisclosure } from '@/hooks/use-disclosure';

import classes from './TextField.module.css';

interface TextFieldProps {
	label: string;
	description?: string;
	required?: boolean;
	error?: string;
	direction?: 'row' | 'column';
}

interface TextWrapperProps extends TextFieldProps {
	inputId: string;
	descriptionId: string;
	errorId: string;
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
				aria-required={required}
				data-error={error ? 'true' : undefined}
				required={required}
				{...props}
			/>
		</TextWrapper>
	);
};

const PasswordInput = ({
	label,
	description,
	error,
	id,
	required,
	disabled,
	...props
}: TextInputProps) => {
	const inputId = id || `input-${label}`;
	const descriptionId = `${inputId}-description`;
	const errorId = `${inputId}-error`;

	const [visible, { toggle }] = useDisclosure();

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
			<div style={{ position: 'relative' }}>
				<input
					id={inputId}
					className={classes['password-input']}
					aria-describedby={
						`${description ? descriptionId : ''} ${error ? errorId : ''}`.trim() ||
						undefined
					}
					aria-invalid={Boolean(error)}
					aria-required={required}
					data-error={error ? 'true' : undefined}
					required={required}
					disabled={disabled}
					type={visible ? 'text' : 'password'}
					{...props}
				/>
				<IconButton
					className={classes['password-button']}
					label="Toggle password visibility"
					disabled={disabled}
					tabIndex={-1}
					onMouseDown={(event) => {
						event.preventDefault();
						toggle();
					}}
					onKeyDown={(event) => {
						if (event.key === ' ') {
							event.preventDefault();
							toggle();
						}
					}}
				>
					{visible ? <EyeNoneIcon /> : <EyeOpenIcon />}
				</IconButton>
			</div>
		</TextWrapper>
	);
};

const TextArea = ({
	label,
	description,
	error,
	id,
	required,
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
			required={required}
		>
			<textarea
				id={textareaId}
				aria-describedby={
					`${description ? descriptionId : ''} ${error ? errorId : ''}`.trim() ||
					undefined
				}
				aria-invalid={Boolean(error)}
				aria-required={required}
				{...props}
			/>
		</TextWrapper>
	);
};

export { TextArea, TextInput, PasswordInput };
