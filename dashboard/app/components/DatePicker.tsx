import * as Dialog from '@radix-ui/react-dialog';
import * as VisuallyHidden from '@radix-ui/react-visually-hidden';
import { useState } from 'react';
import { DayPicker, type DateRange } from 'react-day-picker';

import { Group } from '@/components/layout/Flex';

import { CloseButton } from './Button';

import dayPickerClasses from 'react-day-picker/style.module.css';
import classes from './DatePicker.module.css';

interface DatePickerProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

const DatePickerRange = ({ open, setOpen }: DatePickerProps) => {
	const [date, setDate] = useState<DateRange | undefined>();

	return (
		<Dialog.Root open={open} onOpenChange={setOpen} defaultOpen>
			<Dialog.Portal>
				<Dialog.Overlay className={classes.overlay} />
				<Dialog.Content className={classes.content}>
					<Group className={classes.header}>
						<Dialog.Title className={classes.title}>
							Select a date range
						</Dialog.Title>
						<Dialog.Close asChild>
							<CloseButton label="Close date picker" />
						</Dialog.Close>
					</Group>
					<VisuallyHidden.Root asChild>
						<Dialog.Description>
							Select the start and end date for the date range.
						</Dialog.Description>
					</VisuallyHidden.Root>

					<DayPicker
						mode="range"
						classNames={{
							...dayPickerClasses,
							months: undefined,
						}}
						selected={date}
						onSelect={setDate}
						numberOfMonths={2}
					/>
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>
	);
};

export { classes as datePickerClasses, DatePickerRange };
