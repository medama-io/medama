import * as CheckboxPrimitive from '@radix-ui/react-checkbox';
import { CheckIcon } from '@radix-ui/react-icons';
import * as Label from '@radix-ui/react-label';
import type React from 'react';

import { InfoTooltip } from '@/components/InfoTooltip';
import { Group } from '@/components/layout/Flex';

import classes from './Checkbox.module.css';

interface CheckBoxProps {
	label: string;
	value: string;
	icon?: React.ReactNode;
	tooltip?: React.ReactNode;

	checked?: boolean;
	disabled?: boolean;
	onCheckedChange?: (checked: boolean) => void;
}

const CheckBox = ({
	label,
	value,
	icon,
	tooltip,
	checked,
	disabled,
	onCheckedChange,
}: CheckBoxProps) => {
	const id = `checkbox-${value}`;
	return (
		<Group style={{ justifyContent: 'flex-start' }}>
			<CheckboxPrimitive.Root
				id={id}
				className={classes.root}
				value={value}
				checked={checked}
				disabled={disabled}
				onCheckedChange={onCheckedChange}
			>
				<CheckboxPrimitive.Indicator className={classes.indicator}>
					{icon ? icon : <CheckIcon />}
				</CheckboxPrimitive.Indicator>
			</CheckboxPrimitive.Root>
			<Group>
				<Label.Root className={classes.label} htmlFor={id}>
					{label}
				</Label.Root>
				{tooltip && <InfoTooltip>{tooltip}</InfoTooltip>}
			</Group>
		</Group>
	);
};

export { CheckBox };
