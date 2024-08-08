import * as CheckboxPrimitive from '@radix-ui/react-checkbox';
import { CheckIcon } from '@radix-ui/react-icons';

import { InfoTooltip } from '@/components/InfoTooltip';
import { Group } from '@/components/layout/Flex';

import classes from './Checkbox.module.css';

interface CheckBoxProps {
	label: string;
	value: string;
	icon?: React.ReactNode;

	checked?: boolean;
	disabled?: boolean;
	onCheckedChange?: (checked: boolean) => void;
}

const CheckBox = ({
	label,
	value,
	icon,
	checked,
	disabled,
	onCheckedChange,
}: CheckBoxProps) => (
	<Group style={{ justifyContent: 'flex-start' }}>
		<CheckboxPrimitive.Root
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
		<label className={classes.label} htmlFor="c1">
			{label}
		</label>
		<InfoTooltip />
	</Group>
);

export { CheckBox };
