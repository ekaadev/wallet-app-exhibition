<script lang="ts">
    import {Button} from '$lib/components/ui/button/index.js';
    import {Input} from '$lib/components/ui/input/index.js';
    import {Field, FieldDescription, FieldGroup, FieldLabel} from '$lib/components/ui/field/index.js';
    import * as Card from '$lib/components/ui/card/index.js';

    // Generate a unique ID for form fields
    let id = Math.random().toString(36).substring(2, 9);
</script>

<Card.Root class="mx-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">Login</Card.Title>
        <Card.Description>Enter your email below to login to your account</Card.Description>
    </Card.Header>
    <Card.Content>
        <form>
            <FieldGroup>
                <Field>
                    <FieldLabel for="email-{id}">Email</FieldLabel>
                    <Input id="email-{id}" type="email" placeholder="m@example.com" required/>
                </Field>
                <Field>
                    <div class="flex items-center">
                        <FieldLabel for="password-{id}">Password</FieldLabel>
                    </div>
                    <Input id="password-{id}" type="password" placeholder="********" required/>
                </Field>
                <Field>
                    <Button type="submit" class="w-full">Login</Button>
                    <FieldDescription class="text-center">
                        Don't have an account? <a href="/signup">Sign up</a>
                    </FieldDescription>
                </Field>
            </FieldGroup>
        </form>
    </Card.Content>
</Card.Root>
