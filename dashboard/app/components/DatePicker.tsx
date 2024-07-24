import * as Dialog from '@radix-ui/react-dialog';
import * as VisuallyHidden from '@radix-ui/react-visually-hidden';
import { sub } from 'date-fns';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useState } from 'react';
import { DayPicker, type DateRange } from 'react-day-picker';

import { Button, CloseButton } from '@/components/Button';
import { Group } from '@/components/layout/Flex';

import dayPickerClasses from 'react-day-picker/style.module.css';
import classes from './DatePicker.module.css';

interface DatePickerProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

const DatePickerRange = ({ open, setOpen }: DatePickerProps) => {
	const [date, setDate] = useState<DateRange | undefined>();

	return (
		<>
			<Dialog.Root open={open} onOpenChange={setOpen}>
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
								root: classes.root,
								months: classes.months,
								chevron: undefined,
								month_grid: classes.month_grid,
							}}
							selected={date}
							onSelect={setDate}
							numberOfMonths={2}
							components={{
								Chevron: (props) => {
									if (props.orientation === 'left') {
										return <ChevronLeft {...props} />;
									}
									return <ChevronRight {...props} />;
								},
							}}
							defaultMonth={sub(new Date(), { months: 1 })}
							startMonth={new Date(2024, 0)}
							endMonth={new Date()}
						/>
						<Button className={classes.apply} onClick={() => setOpen(false)}>
							Apply
						</Button>
					</Dialog.Content>
				</Dialog.Portal>
			</Dialog.Root>
		</>
	);
};

export { classes as datePickerClasses, DatePickerRange };
