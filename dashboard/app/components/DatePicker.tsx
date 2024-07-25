import * as Dialog from '@radix-ui/react-dialog';
import * as VisuallyHidden from '@radix-ui/react-visually-hidden';
import { useSearchParams } from '@remix-run/react';
import { formatISO, parseISO, sub } from 'date-fns';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useState } from 'react';
import { type DateRange, DayPicker } from 'react-day-picker';

import { Button, CloseButton } from '@/components/Button';
import { Group } from '@/components/layout/Flex';

import { useMediaQuery } from '@mantine/hooks';
import classes from './DatePicker.module.css';

interface DatePickerProps {
	open: boolean;
	setOpen: (open: boolean) => void;
}

const DatePickerRange = ({ open, setOpen }: DatePickerProps) => {
	const [searchParams, setSearchParams] = useSearchParams();
	const [date, setDate] = useState<DateRange | undefined>({
		from: searchParams.get('start')
			? // biome-ignore lint/style/noNonNullAssertion: Valid.
				parseISO(searchParams.get('start')!)
			: undefined,

		to: searchParams.get('end')
			? // biome-ignore lint/style/noNonNullAssertion: Valid.
				parseISO(searchParams.get('end')!)
			: undefined,
	});

	const handleSubmit = () => {
		if (date) {
			setSearchParams((params) => {
				params.set('period', 'custom');
				if (date.from)
					params.set('start', formatISO(date.from, { representation: 'date' }));
				if (date.to)
					params.set('end', formatISO(date.to, { representation: 'date' }));
				return params;
			});
		}
		setOpen(false);
	};

	const isMobile = useMediaQuery('(max-width: 48em)');

	return (
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
						classNames={classes}
						selected={date}
						onSelect={setDate}
						numberOfMonths={isMobile ? 1 : 2}
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
					<Button className={classes.apply} onClick={handleSubmit}>
						Apply
					</Button>
				</Dialog.Content>
			</Dialog.Portal>
		</Dialog.Root>
	);
};

export { classes as datePickerClasses, DatePickerRange };
