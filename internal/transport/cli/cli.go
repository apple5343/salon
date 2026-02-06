package cli

import (
	"salon/internal/models"
	"salon/internal/service"

	"github.com/spf13/cobra"
)

type CLI struct {
	root            *cobra.Command
	employeeService service.EmployeeService
}

func NewCLI(employeeService service.EmployeeService) *CLI {
	return &CLI{
		root: &cobra.Command{
			Use:   "app",
			Short: "application cli",
		},
		employeeService: employeeService,
	}
}

func (c *CLI) RegisterAdmin() *cobra.Command {
	var admin models.Employee
	run := func(cmd *cobra.Command, args []string) {
		admin.Role = models.EmployeeRoleAdmin
		admin.Status = models.EmployeeStatusActive
		a, err := c.employeeService.RegisterAdmin(cmd.Context(), &admin)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		cmd.Printf("admin registered: \n")
		cmd.Printf("\tID: %s\n", a.ID)
		cmd.Printf("\tFull name: %s\n", a.FullName)
		cmd.Printf("\tEmail: %s\n", a.Email)
		cmd.Printf("\tPhone: %s\n", a.Phone)
		cmd.Printf("\tPassport series: %s\n", a.Passport.Series)
		cmd.Printf("\tPassport number: %s\n", a.Passport.Number)
		cmd.Printf("\tPassport issued by: %s\n", a.Passport.IssuedBy)
		cmd.Printf("\tHire date: %s\n", a.HireDate)
	}
	registerCmd := &cobra.Command{
		Use:   "register-admin",
		Short: "register admin",
		Run:   run,
	}
	registerCmd.Flags().StringVar(&admin.FullName, "fullname", "", "Имя админа")
	registerCmd.Flags().StringVar(&admin.Phone, "phone", "", "Телефон админа")
	registerCmd.Flags().StringVar(&admin.Email, "email", "", "Почта админа")
	registerCmd.Flags().StringVar(&admin.PasswordHash, "password", "", "Пароль админа")
	registerCmd.Flags().StringVar(&admin.Passport.Series, "passport_series", "", "Cерия паспорта админа")
	registerCmd.Flags().StringVar(&admin.Passport.Number, "passport_number", "", "Номер паспорта админа")
	registerCmd.Flags().StringVar(&admin.Passport.IssuedBy, "passport_issued_by", "", "Кем выдан паспорт админа")
	registerCmd.MarkFlagsRequiredTogether("fullname", "phone", "email", "password", "passport_series", "passport_number", "passport_issued_by")
	return registerCmd
}

func (c *CLI) Run() error {
	c.root.AddCommand(c.RegisterAdmin())
	if err := c.root.Execute(); err != nil {
		return err
	}
	return nil
}
