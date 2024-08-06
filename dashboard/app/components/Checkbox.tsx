import * as CheckboxPrimitive from '@radix-ui/react-checkbox';
import { Check } from 'lucide-react';

import { Group } from '@/components/layout/Flex';
import { InfoTooltip } from '@/components/InfoTooltip';

import classes from './Checkbox.module.css';

interface CheckBoxProps {
	label: string;
	icon?: React.ReactNode;

	checked?: boolean;
	disabled?: boolean;
}

const CheckBox = ({ label, icon, checked, disabled }: CheckBoxProps) => (
	<Group style={{ justifyContent: 'flex-start' }}>
		<CheckboxPrimitive.Root
			className={classes.root}
			checked={checked}
			disabled={disabled}
		>
			<CheckboxPrimitive.Indicator className={classes.indicator}>
				{icon ? icon : <Check />}
			</CheckboxPrimitive.Indicator>
		</CheckboxPrimitive.Root>
		<label className={classes.label} htmlFor="c1">
			{label}
		</label>
		<InfoTooltip />
	</Group>
);

export { CheckBox };
