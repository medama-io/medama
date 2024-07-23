import { useState } from 'react';
import { DayPicker, type DateRange } from 'react-day-picker';
import * as Dialog from '@radix-ui/react-dialog';

import dayPickerClasses from 'react-day-picker/style.module.css';
import classes from './DatePicker.module.css';

interface DatePickerProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

const DatePickerRange = ({ open, setOpen }: DatePickerProps) => {
	const [date, setDate] = useState<DateRange | undefined>();

	return (
		<Dialog.Root open={open} onOpenChange={setOpen}>
			<Dialog.Portal>
				<Dialog.Overlay className={classes.overlay} />
				<Dialog.Content className={classes.content}>
					<Dialog.Title />
					<Dialog.Description />
					<Dialog.Close />
					<DayPicker
						mode="range"
						classNames={dayPickerClasses}
						selected={date}
						onSelect={setDate}
					/>
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>
	);
};

export { DatePickerRange, classes as datePickerClasses };
